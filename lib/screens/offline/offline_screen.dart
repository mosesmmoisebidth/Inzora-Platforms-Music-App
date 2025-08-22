import 'package:flutter/material.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter_staggered_animations/flutter_staggered_animations.dart';
import '../../models/song.dart';
import '../../services/audio_service.dart';
import '../../services/storage_service.dart';

class OfflineScreen extends StatefulWidget {
  const OfflineScreen({super.key});

  @override
  State<OfflineScreen> createState() => _OfflineScreenState();
}

class _OfflineScreenState extends State<OfflineScreen> {
  final AudioService _audioService = AudioService();
  final StorageService _storageService = StorageService();

  List<Song> _offlineSongs = [];
  bool _isLoading = true;
  Song? _currentSong;
  bool _isSelectionMode = false;
  Set<String> _selectedSongs = {};

  @override
  void initState() {
    super.initState();
    _loadOfflineSongs();
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

    setState(() {
      _currentSong = _audioService.currentSong;
    });
  }

  Future<void> _loadOfflineSongs() async {
    setState(() {
      _isLoading = true;
    });

    try {
      final offlineSongs = _storageService.getOfflineSongs();
      if (mounted) {
        setState(() {
          _offlineSongs = offlineSongs;
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
      await _audioService.playFromPlaylist(_offlineSongs, index);
      await _storageService.addToRecentlyPlayed(song);
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error playing song: $e')),
      );
    }
  }

  Future<void> _removeFromOffline(Song song) async {
    try {
      await _storageService.removeFromOffline(song);
      setState(() {
        _offlineSongs.removeWhere((s) => s.id == song.id);
      });
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Removed from offline downloads')),
      );
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error removing from offline: $e')),
      );
    }
  }

  Future<void> _removeSelectedSongs() async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        backgroundColor: const Color(0xFF1A1A2E),
        title: const Text('Remove Downloads', style: TextStyle(color: Colors.white)),
        content: Text(
          'Are you sure you want to remove ${_selectedSongs.length} song${_selectedSongs.length > 1 ? 's' : ''} from offline downloads?',
          style: const TextStyle(color: Colors.white70),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.of(context).pop(true),
            style: ElevatedButton.styleFrom(backgroundColor: Colors.red),
            child: const Text('Remove'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      try {
        for (final songId in _selectedSongs) {
          final song = _offlineSongs.firstWhere((s) => s.id == songId);
          await _storageService.removeFromOffline(song);
        }
        
        setState(() {
          _offlineSongs.removeWhere((song) => _selectedSongs.contains(song.id));
          _selectedSongs.clear();
          _isSelectionMode = false;
        });

        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Selected songs removed from downloads')),
        );
      } catch (e) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error removing songs: $e')),
        );
      }
    }
  }

  Future<void> _clearAllOffline() async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        backgroundColor: const Color(0xFF1A1A2E),
        title: const Text('Clear All Downloads', style: TextStyle(color: Colors.white)),
        content: const Text(
          'Are you sure you want to remove all offline downloads? This action cannot be undone.',
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
            child: const Text('Clear All'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      try {
        for (final song in _offlineSongs) {
          await _storageService.removeFromOffline(song);
        }
        setState(() {
          _offlineSongs.clear();
        });
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('All offline downloads cleared')),
        );
      } catch (e) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error clearing downloads: $e')),
        );
      }
    }
  }

  void _toggleSelectionMode() {
    setState(() {
      _isSelectionMode = !_isSelectionMode;
      _selectedSongs.clear();
    });
  }

  void _toggleSongSelection(String songId) {
    setState(() {
      if (_selectedSongs.contains(songId)) {
        _selectedSongs.remove(songId);
      } else {
        _selectedSongs.add(songId);
      }
    });
  }

  void _selectAllSongs() {
    setState(() {
      _selectedSongs = Set.from(_offlineSongs.map((song) => song.id));
    });
  }

  Future<void> _playAll() async {
    if (_offlineSongs.isEmpty) return;
    
    try {
      await _audioService.setPlaylist(_offlineSongs);
      await _audioService.play();
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error playing songs: $e')),
      );
    }
  }

  double _getTotalSize() {
    // Simulate download size calculation (in MB)
    // In a real app, this would come from actual file sizes
    return _offlineSongs.length * 3.5; // Assuming average 3.5MB per song
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
              if (_offlineSongs.isNotEmpty) _buildStorageInfo(),
              Expanded(
                child: _isLoading
                    ? const Center(child: CircularProgressIndicator())
                    : _offlineSongs.isEmpty
                        ? _buildEmptyState()
                        : _buildOfflineSongsList(),
              ),
            ],
          ),
        ),
      ),
      bottomNavigationBar: _isSelectionMode && _selectedSongs.isNotEmpty
          ? _buildSelectionBottomBar()
          : null,
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
                onPressed: () {
                  if (_isSelectionMode) {
                    _toggleSelectionMode();
                  } else {
                    Navigator.of(context).pop();
                  }
                },
                icon: Icon(
                  _isSelectionMode ? Icons.close : Icons.arrow_back,
                  color: Colors.white,
                ),
              ),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      _isSelectionMode 
                          ? '${_selectedSongs.length} selected'
                          : 'Downloaded Music',
                      style: const TextStyle(
                        fontSize: 28,
                        fontWeight: FontWeight.bold,
                        color: Colors.white,
                      ),
                    ),
                    if (!_isSelectionMode)
                      Text(
                        '${_offlineSongs.length} songs • ${_getTotalSize().toStringAsFixed(1)} MB',
                        style: TextStyle(
                          fontSize: 16,
                          color: Colors.white.withOpacity(0.7),
                        ),
                      ),
                  ],
                ),
              ),
              if (!_isSelectionMode && _offlineSongs.isNotEmpty) ...[
                IconButton(
                  onPressed: _toggleSelectionMode,
                  icon: const Icon(Icons.checklist, color: Colors.white),
                ),
                PopupMenuButton<String>(
                  icon: const Icon(Icons.more_vert, color: Colors.white),
                  color: const Color(0xFF1A1A2E),
                  onSelected: (value) {
                    switch (value) {
                      case 'clear_all':
                        _clearAllOffline();
                        break;
                      case 'refresh':
                        _loadOfflineSongs();
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
                      value: 'clear_all',
                      child: Row(
                        children: [
                          Icon(Icons.delete_sweep, color: Colors.red, size: 20),
                          SizedBox(width: 12),
                          Text('Clear All', style: TextStyle(color: Colors.red)),
                        ],
                      ),
                    ),
                  ],
                ),
              ] else if (_isSelectionMode) ...[
                TextButton(
                  onPressed: _selectAllSongs,
                  child: const Text('Select All'),
                ),
              ],
            ],
          ),
          if (_offlineSongs.isNotEmpty && !_isSelectionMode) ...[
            const SizedBox(height: 16),
            SizedBox(
              width: double.infinity,
              child: ElevatedButton.icon(
                onPressed: _playAll,
                icon: const Icon(Icons.play_arrow),
                label: const Text('Play All Downloaded'),
                style: ElevatedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(vertical: 12),
                ),
              ),
            ),
          ],
        ],
      ),
    );
  }

  Widget _buildStorageInfo() {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 20, vertical: 8),
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.05),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: Colors.white.withOpacity(0.1),
          width: 1,
        ),
      ),
      child: Row(
        children: [
          const Icon(
            Icons.storage,
            color: Color(0xFFE94560),
            size: 24,
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Storage Used',
                  style: TextStyle(
                    fontSize: 14,
                    color: Colors.white.withOpacity(0.7),
                  ),
                ),
                Text(
                  '${_getTotalSize().toStringAsFixed(1)} MB of downloaded music',
                  style: const TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                    color: Colors.white,
                  ),
                ),
              ],
            ),
          ),
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
              Icons.offline_pin,
              size: 60,
              color: Color(0xFFE94560),
            ),
          ),
          const SizedBox(height: 24),
          const Text(
            'No Downloaded Music',
            style: TextStyle(
              fontSize: 24,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'Download your favorite songs to listen\nwhen you\'re offline',
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
            child: const Text('Browse Music'),
          ),
        ],
      ),
    );
  }

  Widget _buildOfflineSongsList() {
    return AnimationLimiter(
      child: ListView.builder(
        padding: const EdgeInsets.symmetric(horizontal: 20.0),
        itemCount: _offlineSongs.length,
        itemBuilder: (context, index) {
          return AnimationConfiguration.staggeredList(
            position: index,
            duration: const Duration(milliseconds: 375),
            child: SlideAnimation(
              verticalOffset: 50.0,
              child: FadeInAnimation(
                child: _buildSongTile(_offlineSongs[index], index),
              ),
            ),
          );
        },
      ),
    );
  }

  Widget _buildSongTile(Song song, int index) {
    final isCurrentSong = _currentSong?.id == song.id;
    final isSelected = _selectedSongs.contains(song.id);

    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      decoration: BoxDecoration(
        color: isCurrentSong 
            ? const Color(0xFFE94560).withOpacity(0.1)
            : isSelected
                ? const Color(0xFFE94560).withOpacity(0.2)
                : Colors.white.withOpacity(0.05),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: isCurrentSong 
              ? const Color(0xFFE94560).withOpacity(0.3)
              : isSelected
                  ? const Color(0xFFE94560)
                  : Colors.white.withOpacity(0.1),
          width: 1,
        ),
      ),
      child: ListTile(
        contentPadding: const EdgeInsets.all(12),
        leading: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            if (_isSelectionMode) ...[
              Checkbox(
                value: isSelected,
                onChanged: (value) => _toggleSongSelection(song.id),
                fillColor: MaterialStateProperty.all(const Color(0xFFE94560)),
              ),
              const SizedBox(width: 8),
            ],
            Container(
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
                  Positioned(
                    bottom: 2,
                    right: 2,
                    child: Container(
                      padding: const EdgeInsets.all(2),
                      decoration: BoxDecoration(
                        color: Colors.green,
                        shape: BoxShape.circle,
                      ),
                      child: const Icon(
                        Icons.offline_pin,
                        color: Colors.white,
                        size: 12,
                      ),
                    ),
                  ),
                  if (isCurrentSong && !_isSelectionMode)
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
          ],
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
              '${song.durationText} • Downloaded',
              style: TextStyle(
                fontSize: 12,
                color: Colors.green.withOpacity(0.8),
              ),
            ),
          ],
        ),
        trailing: _isSelectionMode
            ? null
            : Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  IconButton(
                    icon: const Icon(
                      Icons.delete_outline,
                      color: Colors.red,
                      size: 20,
                    ),
                    onPressed: () => _removeFromOffline(song),
                  ),
                  IconButton(
                    icon: const Icon(
                      Icons.play_arrow,
                      color: Color(0xFFE94560),
                      size: 24,
                    ),
                    onPressed: () => _playSong(song, index),
                  ),
                ],
              ),
        onTap: () {
          if (_isSelectionMode) {
            _toggleSongSelection(song.id);
          } else {
            _playSong(song, index);
          }
        },
        onLongPress: () {
          if (!_isSelectionMode) {
            _toggleSelectionMode();
            _toggleSongSelection(song.id);
          }
        },
      ),
    );
  }

  Widget _buildSelectionBottomBar() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: const Color(0xFF1A1A2E),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.3),
            blurRadius: 10,
            offset: const Offset(0, -2),
          ),
        ],
      ),
      child: SafeArea(
        child: Row(
          children: [
            Expanded(
              child: Text(
                '${_selectedSongs.length} song${_selectedSongs.length > 1 ? 's' : ''} selected',
                style: const TextStyle(
                  color: Colors.white,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ),
            ElevatedButton.icon(
              onPressed: _removeSelectedSongs,
              icon: const Icon(Icons.delete, size: 18),
              label: const Text('Remove'),
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.red,
                foregroundColor: Colors.white,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
