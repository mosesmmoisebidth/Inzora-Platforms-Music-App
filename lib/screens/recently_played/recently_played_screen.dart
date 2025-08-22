import 'package:flutter/material.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter_staggered_animations/flutter_staggered_animations.dart';
import '../../models/song.dart';
import '../../services/audio_service.dart';
import '../../services/storage_service.dart';

class RecentlyPlayedScreen extends StatefulWidget {
  const RecentlyPlayedScreen({super.key});

  @override
  State<RecentlyPlayedScreen> createState() => _RecentlyPlayedScreenState();
}

class _RecentlyPlayedScreenState extends State<RecentlyPlayedScreen> {
  final AudioService _audioService = AudioService();
  final StorageService _storageService = StorageService();

  List<Song> _recentlyPlayed = [];
  bool _isLoading = true;
  Song? _currentSong;

  @override
  void initState() {
    super.initState();
    _loadRecentlyPlayed();
    _setupListeners();
  }

  void _setupListeners() {
    _audioService.currentSongStream.listen((song) {
      if (mounted) {
        setState(() {
          _currentSong = song;
        });
      }
    });

    // Initialize with current song
    setState(() {
      _currentSong = _audioService.currentSong;
    });
  }

  Future<void> _loadRecentlyPlayed() async {
    setState(() {
      _isLoading = true;
    });

    try {
      final recentlyPlayed = _storageService.getRecentlyPlayed();
      if (mounted) {
        setState(() {
          _recentlyPlayed = recentlyPlayed;
          _isLoading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  Future<void> _playSong(Song song, int index) async {
    try {
      await _audioService.playFromPlaylist(_recentlyPlayed, index);
      await _storageService.addToRecentlyPlayed(song);
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error playing song: $e')),
      );
    }
  }

  Future<void> _toggleFavorite(Song song) async {
    try {
      final isFavorite = _storageService.isFavorite(song);
      if (isFavorite) {
        await _storageService.removeFromFavorites(song);
      } else {
        await _storageService.addToFavorites(song);
      }
      
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(isFavorite ? 'Removed from favorites' : 'Added to favorites'),
          duration: const Duration(seconds: 1),
        ),
      );
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error updating favorites: $e')),
      );
    }
  }

  Future<void> _clearHistory() async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        backgroundColor: const Color(0xFF1A1A2E),
        title: const Text('Clear History', style: TextStyle(color: Colors.white)),
        content: const Text(
          'Are you sure you want to clear your listening history? This action cannot be undone.',
          style: TextStyle(color: Colors.white70),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.of(context).pop(true),
            style: ElevatedButton.styleFrom(backgroundColor: Colors.red),
            child: const Text('Clear'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      try {
        await _storageService.clearRecentlyPlayed();
        setState(() {
          _recentlyPlayed.clear();
        });
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Listening history cleared')),
        );
      } catch (e) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error clearing history: $e')),
        );
      }
    }
  }

  Future<void> _playAll() async {
    if (_recentlyPlayed.isEmpty) return;
    
    try {
      await _audioService.setPlaylist(_recentlyPlayed);
      await _audioService.play();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error playing songs: $e')),
      );
    }
  }

  Future<void> _shufflePlay() async {
    if (_recentlyPlayed.isEmpty) return;
    
    try {
      await _audioService.setPlaylist(_recentlyPlayed);
      _audioService.toggleShuffle();
      await _audioService.play();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error playing songs: $e')),
      );
    }
  }

  String _getTimeAgo(DateTime dateTime) {
    final now = DateTime.now();
    final difference = now.difference(dateTime);

    if (difference.inDays > 0) {
      return '${difference.inDays} day${difference.inDays > 1 ? 's' : ''} ago';
    } else if (difference.inHours > 0) {
      return '${difference.inHours} hour${difference.inHours > 1 ? 's' : ''} ago';
    } else if (difference.inMinutes > 0) {
      return '${difference.inMinutes} minute${difference.inMinutes > 1 ? 's' : ''} ago';
    } else {
      return 'Just now';
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFF0D1421),
      body: Container(
        decoration: const BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
            colors: [
              Color(0xFF0D1421),
              Color(0xFF1A1A2E),
              Color(0xFF16213E),
            ],
          ),
        ),
        child: SafeArea(
          child: Column(
            children: [
              _buildHeader(),
              Expanded(
                child: _isLoading
                    ? const Center(child: CircularProgressIndicator())
                    : _recentlyPlayed.isEmpty
                        ? _buildEmptyState()
                        : _buildRecentlyPlayedList(),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildHeader() {
    return Padding(
      padding: const EdgeInsets.all(20.0),
      child: Column(
        children: [
          Row(
            children: [
              IconButton(
                onPressed: () => Navigator.of(context).pop(),
                icon: const Icon(Icons.arrow_back, color: Colors.white),
              ),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Recently Played',
                      style: TextStyle(
                        fontSize: 28,
                        fontWeight: FontWeight.bold,
                        color: Colors.white,
                      ),
                    ),
                    Text(
                      '${_recentlyPlayed.length} songs in history',
                      style: TextStyle(
                        fontSize: 16,
                        color: Colors.white.withOpacity(0.7),
                      ),
                    ),
                  ],
                ),
              ),
              if (_recentlyPlayed.isNotEmpty)
                PopupMenuButton<String>(
                  icon: const Icon(Icons.more_vert, color: Colors.white),
                  color: const Color(0xFF1A1A2E),
                  onSelected: (value) {
                    switch (value) {
                      case 'clear':
                        _clearHistory();
                        break;
                      case 'refresh':
                        _loadRecentlyPlayed();
                        break;
                    }
                  },
                  itemBuilder: (context) => [
                    const PopupMenuItem(
                      value: 'refresh',
                      child: Row(
                        children: [
                          Icon(Icons.refresh, color: Colors.white, size: 20),
                          SizedBox(width: 12),
                          Text('Refresh', style: TextStyle(color: Colors.white)),
                        ],
                      ),
                    ),
                    const PopupMenuItem(
                      value: 'clear',
                      child: Row(
                        children: [
                          Icon(Icons.clear_all, color: Colors.red, size: 20),
                          SizedBox(width: 12),
                          Text('Clear History', style: TextStyle(color: Colors.red)),
                        ],
                      ),
                    ),
                  ],
                ),
            ],
          ),
          if (_recentlyPlayed.isNotEmpty) ...[
            const SizedBox(height: 20),
            Row(
              children: [
                Expanded(
                  child: ElevatedButton.icon(
                    onPressed: _playAll,
                    icon: const Icon(Icons.play_arrow),
                    label: const Text('Play All'),
                    style: ElevatedButton.styleFrom(
                      padding: const EdgeInsets.symmetric(vertical: 12),
                    ),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: OutlinedButton.icon(
                    onPressed: _shufflePlay,
                    icon: const Icon(Icons.shuffle),
                    label: const Text('Shuffle'),
                    style: OutlinedButton.styleFrom(
                      foregroundColor: Colors.white,
                      side: BorderSide(color: Colors.white.withOpacity(0.3)),
                      padding: const EdgeInsets.symmetric(vertical: 12),
                    ),
                  ),
                ),
              ],
            ),
          ],
        ],
      ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            width: 120,
            height: 120,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              color: const Color(0xFFE94560).withOpacity(0.1),
              border: Border.all(
                color: const Color(0xFFE94560).withOpacity(0.3),
                width: 2,
              ),
            ),
            child: const Icon(
              Icons.history,
              size: 60,
              color: Color(0xFFE94560),
            ),
          ),
          const SizedBox(height: 24),
          const Text(
            'No Recently Played Songs',
            style: TextStyle(
              fontSize: 24,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'Your listening history will appear here\nonce you start playing music',
            textAlign: TextAlign.center,
            style: TextStyle(
              fontSize: 16,
              color: Colors.white.withOpacity(0.6),
              height: 1.5,
            ),
          ),
          const SizedBox(height: 32),
          ElevatedButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Discover Music'),
          ),
        ],
      ),
    );
  }

  Widget _buildRecentlyPlayedList() {
    return AnimationLimiter(
      child: ListView.builder(
        padding: const EdgeInsets.symmetric(horizontal: 20.0),
        itemCount: _recentlyPlayed.length,
        itemBuilder: (context, index) {
          return AnimationConfiguration.staggeredList(
            position: index,
            duration: const Duration(milliseconds: 375),
            child: SlideAnimation(
              verticalOffset: 50.0,
              child: FadeInAnimation(
                child: _buildSongTile(_recentlyPlayed[index], index),
              ),
            ),
          );
        },
      ),
    );
  }

  Widget _buildSongTile(Song song, int index) {
    final isCurrentSong = _currentSong?.id == song.id;
    final isFavorite = _storageService.isFavorite(song);
    
    // Simulate played time (in real app, this would come from the song data)
    final playedTime = DateTime.now().subtract(Duration(hours: index, minutes: index * 15));

    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      decoration: BoxDecoration(
        color: isCurrentSong 
            ? const Color(0xFFE94560).withOpacity(0.1)
            : Colors.white.withOpacity(0.05),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: isCurrentSong 
              ? const Color(0xFFE94560).withOpacity(0.3)
              : Colors.white.withOpacity(0.1),
          width: 1,
        ),
      ),
      child: ListTile(
        contentPadding: const EdgeInsets.all(12),
        leading: Container(
          width: 50,
          height: 50,
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(8),
          ),
          child: Stack(
            children: [
              ClipRRect(
                borderRadius: BorderRadius.circular(8),
                child: CachedNetworkImage(
                  imageUrl: song.albumArt,
                  fit: BoxFit.cover,
                  placeholder: (context, url) => Container(
                    color: Colors.grey[800],
                    child: const Icon(Icons.music_note, size: 24, color: Colors.white54),
                  ),
                  errorWidget: (context, url, error) => Container(
                    color: Colors.grey[800],
                    child: const Icon(Icons.music_note, size: 24, color: Colors.white54),
                  ),
                ),
              ),
              if (isCurrentSong)
                Container(
                  decoration: BoxDecoration(
                    color: Colors.black.withOpacity(0.5),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: const Center(
                    child: Icon(
                      Icons.equalizer,
                      color: Color(0xFFE94560),
                      size: 20,
                    ),
                  ),
                ),
            ],
          ),
        ),
        title: Text(
          song.title,
          style: TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.w600,
            color: isCurrentSong ? const Color(0xFFE94560) : Colors.white,
          ),
          maxLines: 1,
          overflow: TextOverflow.ellipsis,
        ),
        subtitle: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              song.artist,
              style: TextStyle(
                fontSize: 14,
                color: Colors.white.withOpacity(0.7),
              ),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
            const SizedBox(height: 2),
            Text(
              _getTimeAgo(playedTime),
              style: TextStyle(
                fontSize: 12,
                color: Colors.white.withOpacity(0.5),
              ),
            ),
          ],
        ),
        trailing: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            IconButton(
              icon: Icon(
                isFavorite ? Icons.favorite : Icons.favorite_outline,
                color: isFavorite ? const Color(0xFFE94560) : Colors.white.withOpacity(0.6),
                size: 20,
              ),
              onPressed: () => _toggleFavorite(song),
            ),
            PopupMenuButton<String>(
              icon: Icon(
                Icons.more_vert,
                color: Colors.white.withOpacity(0.6),
                size: 20,
              ),
              color: const Color(0xFF1A1A2E),
              onSelected: (value) {
                switch (value) {
                  case 'add_to_queue':
                    _audioService.addToQueue(song);
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text('Added to queue')),
                    );
                    break;
                  case 'add_to_playlist':
                    // Show playlist selection
                    break;
                  case 'go_to_album':
                    // Navigate to album
                    break;
                }
              },
              itemBuilder: (context) => [
                const PopupMenuItem(
                  value: 'add_to_queue',
                  child: Row(
                    children: [
                      Icon(Icons.queue, color: Colors.white, size: 18),
                      SizedBox(width: 12),
                      Text('Add to queue', style: TextStyle(color: Colors.white, fontSize: 14)),
                    ],
                  ),
                ),
                const PopupMenuItem(
                  value: 'add_to_playlist',
                  child: Row(
                    children: [
                      Icon(Icons.playlist_add, color: Colors.white, size: 18),
                      SizedBox(width: 12),
                      Text('Add to playlist', style: TextStyle(color: Colors.white, fontSize: 14)),
                    ],
                  ),
                ),
                const PopupMenuItem(
                  value: 'go_to_album',
                  child: Row(
                    children: [
                      Icon(Icons.album, color: Colors.white, size: 18),
                      SizedBox(width: 12),
                      Text('Go to album', style: TextStyle(color: Colors.white, fontSize: 14)),
                    ],
                  ),
                ),
              ],
            ),
          ],
        ),
        onTap: () => _playSong(song, index),
      ),
    );
  }
}
