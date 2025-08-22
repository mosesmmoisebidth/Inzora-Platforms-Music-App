import 'package:flutter/material.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter_staggered_animations/flutter_staggered_animations.dart';
import '../../models/playlist.dart';
import '../../models/song.dart';
import '../../services/audio_service.dart';
import '../../services/storage_service.dart';

class PlaylistDetailScreen extends StatefulWidget {
  final Playlist playlist;

  const PlaylistDetailScreen({super.key, required this.playlist});

  @override
  State<PlaylistDetailScreen> createState() => _PlaylistDetailScreenState();
}

class _PlaylistDetailScreenState extends State<PlaylistDetailScreen> {
  final AudioService _audioService = AudioService();
  final StorageService _storageService = StorageService();

  late Playlist _playlist;
  bool _isLoading = false;
  Song? _currentSong;

  @override
  void initState() {
    super.initState();
    _playlist = widget.playlist;
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

  Future<void> _playSong(Song song, int index) async {
    try {
      setState(() {
        _isLoading = true;
      });

      await _audioService.playFromPlaylist(_playlist.songs, index);
      await _storageService.addToRecentlyPlayed(song);
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error playing song: $e')),
      );
    } finally {
      setState(() {
        _isLoading = false;
      });
    }
  }

  Future<void> _playPlaylist() async {
    if (_playlist.songs.isEmpty) return;
    
    try {
      setState(() {
        _isLoading = true;
      });

      await _audioService.setPlaylist(_playlist.songs);
      await _audioService.play();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error playing playlist: $e')),
      );
    } finally {
      setState(() {
        _isLoading = false;
      });
    }
  }

  Future<void> _shufflePlaylist() async {
    if (_playlist.songs.isEmpty) return;
    
    try {
      setState(() {
        _isLoading = true;
      });

      await _audioService.setPlaylist(_playlist.songs);
      _audioService.toggleShuffle();
      await _audioService.play();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error shuffling playlist: $e')),
      );
    } finally {
      setState(() {
        _isLoading = false;
      });
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

  void _showPlaylistOptions() {
    showModalBottomSheet(
      context: context,
      backgroundColor: const Color(0xFF1A1A2E),
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => Container(
        padding: const EdgeInsets.all(20),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 40,
              height: 4,
              decoration: BoxDecoration(
                color: Colors.white.withOpacity(0.3),
                borderRadius: BorderRadius.circular(2),
              ),
            ),
            const SizedBox(height: 20),
            _buildBottomSheetItem(
              icon: Icons.share,
              title: 'Share Playlist',
              onTap: () => Navigator.of(context).pop(),
            ),
            _buildBottomSheetItem(
              icon: Icons.edit,
              title: 'Edit Playlist',
              onTap: () => Navigator.of(context).pop(),
            ),
            _buildBottomSheetItem(
              icon: Icons.download,
              title: 'Download Playlist',
              onTap: () => Navigator.of(context).pop(),
            ),
            _buildBottomSheetItem(
              icon: Icons.info,
              title: 'Playlist Info',
              onTap: () => Navigator.of(context).pop(),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildBottomSheetItem({
    required IconData icon,
    required String title,
    required VoidCallback onTap,
  }) {
    return ListTile(
      leading: Icon(icon, color: Colors.white),
      title: Text(title, style: const TextStyle(color: Colors.white)),
      onTap: onTap,
      contentPadding: EdgeInsets.zero,
    );
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
        child: CustomScrollView(
          slivers: [
            _buildSliverAppBar(),
            SliverToBoxAdapter(child: _buildPlaylistInfo()),
            SliverToBoxAdapter(child: _buildActionButtons()),
            _buildSongsList(),
            const SliverToBoxAdapter(child: SizedBox(height: 100)),
          ],
        ),
      ),
    );
  }

  Widget _buildSliverAppBar() {
    return SliverAppBar(
      expandedHeight: 200.0,
      floating: false,
      pinned: true,
      backgroundColor: const Color(0xFF0D1421),
      flexibleSpace: FlexibleSpaceBar(
        background: Container(
          decoration: BoxDecoration(
            gradient: LinearGradient(
              begin: Alignment.topCenter,
              end: Alignment.bottomCenter,
              colors: [
                Colors.black.withOpacity(0.3),
                Colors.transparent,
              ],
            ),
          ),
          child: CachedNetworkImage(
            imageUrl: _playlist.coverImage,
            fit: BoxFit.cover,
            placeholder: (context, url) => Container(
              decoration: const BoxDecoration(
                gradient: LinearGradient(
                  colors: [
                    Color(0xFFE94560),
                    Color(0xFFF16E00),
                  ],
                ),
              ),
              child: const Icon(
                Icons.playlist_play,
                size: 80,
                color: Colors.white,
              ),
            ),
            errorWidget: (context, url, error) => Container(
              decoration: const BoxDecoration(
                gradient: LinearGradient(
                  colors: [
                    Color(0xFFE94560),
                    Color(0xFFF16E00),
                  ],
                ),
              ),
              child: const Icon(
                Icons.playlist_play,
                size: 80,
                color: Colors.white,
              ),
            ),
          ),
        ),
      ),
      actions: [
        IconButton(
          icon: const Icon(Icons.more_vert, color: Colors.white),
          onPressed: _showPlaylistOptions,
        ),
      ],
    );
  }

  Widget _buildPlaylistInfo() {
    return Padding(
      padding: const EdgeInsets.all(20.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            _playlist.name,
            style: const TextStyle(
              fontSize: 28,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
          const SizedBox(height: 8),
          if (_playlist.description.isNotEmpty) ...[
            Text(
              _playlist.description,
              style: TextStyle(
                fontSize: 16,
                color: Colors.white.withOpacity(0.7),
              ),
            ),
            const SizedBox(height: 8),
          ],
          Row(
            children: [
              Text(
                '${_playlist.songs.length} songs',
                style: TextStyle(
                  fontSize: 14,
                  color: Colors.white.withOpacity(0.6),
                ),
              ),
              if (_playlist.totalDuration > 0) ...[
                Text(
                  ' • ${_playlist.totalDurationText}',
                  style: TextStyle(
                    fontSize: 14,
                    color: Colors.white.withOpacity(0.6),
                  ),
                ),
              ],
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildActionButtons() {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 20.0),
      child: Row(
        children: [
          Expanded(
            child: ElevatedButton.icon(
              onPressed: _isLoading ? null : _playPlaylist,
              icon: _isLoading
                  ? const SizedBox(
                      width: 16,
                      height: 16,
                      child: CircularProgressIndicator(strokeWidth: 2),
                    )
                  : const Icon(Icons.play_arrow),
              label: const Text('Play'),
              style: ElevatedButton.styleFrom(
                padding: const EdgeInsets.symmetric(vertical: 12),
              ),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: OutlinedButton.icon(
              onPressed: _isLoading ? null : _shufflePlaylist,
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
    );
  }

  Widget _buildSongsList() {
    if (_playlist.songs.isEmpty) {
      return SliverToBoxAdapter(
        child: Container(
          padding: const EdgeInsets.all(40.0),
          child: Column(
            children: [
              Icon(
                Icons.music_off,
                size: 64,
                color: Colors.white.withOpacity(0.3),
              ),
              const SizedBox(height: 16),
              Text(
                'No songs in this playlist',
                style: TextStyle(
                  fontSize: 18,
                  color: Colors.white.withOpacity(0.6),
                ),
              ),
            ],
          ),
        ),
      );
    }

    return SliverPadding(
      padding: const EdgeInsets.symmetric(horizontal: 20.0),
      sliver: SliverList(
        delegate: SliverChildBuilderDelegate(
          (context, index) {
            return AnimationConfiguration.staggeredList(
              position: index,
              duration: const Duration(milliseconds: 375),
              child: SlideAnimation(
                verticalOffset: 50.0,
                child: FadeInAnimation(
                  child: _buildSongTile(_playlist.songs[index], index),
                ),
              ),
            );
          },
          childCount: _playlist.songs.length,
        ),
      ),
    );
  }

  Widget _buildSongTile(Song song, int index) {
    final isCurrentSong = _currentSong?.id == song.id;
    final isFavorite = _storageService.isFavorite(song);

    return Container(
      margin: const EdgeInsets.only(bottom: 8),
      decoration: BoxDecoration(
        color: isCurrentSong 
            ? const Color(0xFFE94560).withOpacity(0.1)
            : Colors.white.withOpacity(0.03),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: isCurrentSong 
              ? const Color(0xFFE94560).withOpacity(0.3)
              : Colors.white.withOpacity(0.05),
          width: 1,
        ),
      ),
      child: ListTile(
        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
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
                  width: 50,
                  height: 50,
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
            Text(
              song.durationText,
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
                    break;
                  case 'add_to_playlist':
                    // Show playlist selection
                    break;
                  case 'remove':
                    // Remove from playlist
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
                  value: 'remove',
                  child: Row(
                    children: [
                      Icon(Icons.remove_circle_outline, color: Colors.red, size: 18),
                      SizedBox(width: 12),
                      Text('Remove from playlist', style: TextStyle(color: Colors.red, fontSize: 14)),
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
