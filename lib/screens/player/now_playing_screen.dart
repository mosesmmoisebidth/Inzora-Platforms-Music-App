import 'package:flutter/material.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'dart:math' as math;
import '../../models/song.dart';
import '../../services/audio_service.dart';
import '../../services/storage_service.dart';

class NowPlayingScreen extends StatefulWidget {
  const NowPlayingScreen({super.key});

  @override
  State<NowPlayingScreen> createState() => _NowPlayingScreenState();
}

class _NowPlayingScreenState extends State<NowPlayingScreen>
    with TickerProviderStateMixin {
  final AudioService _audioService = AudioService();
  final StorageService _storageService = StorageService();

  late AnimationController _albumArtController;
  late AnimationController _slideController;
  late Animation<double> _albumArtAnimation;
  late Animation<Offset> _slideAnimation;

  Song? _currentSong;
  PlaybackState _playbackState = PlaybackState.stopped;
  Duration _position = Duration.zero;
  Duration _duration = Duration.zero;
  bool _isDragging = false;
  bool _isFavorite = false;

  @override
  void initState() {
    super.initState();
    _setupAnimations();
    _setupAudioListeners();
  }

  void _setupAnimations() {
    _albumArtController = AnimationController(
      duration: const Duration(seconds: 20),
      vsync: this,
    );

    _slideController = AnimationController(
      duration: const Duration(milliseconds: 300),
      vsync: this,
    );

    _albumArtAnimation = Tween<double>(
      begin: 0,
      end: 2 * math.pi,
    ).animate(_albumArtController);

    _slideAnimation = Tween<Offset>(
      begin: const Offset(0, 1),
      end: Offset.zero,
    ).animate(CurvedAnimation(
      parent: _slideController,
      curve: Curves.easeOutCubic,
    ));

    _slideController.forward();

    // Start rotation if playing
    if (_playbackState == PlaybackState.playing) {
      _albumArtController.repeat();
    }
  }

  void _setupAudioListeners() {
    _audioService.currentSongStream.listen((song) {
      if (mounted) {
        setState(() {
          _currentSong = song;
          _isFavorite = song != null ? _storageService.isFavorite(song) : false;
        });
      }
    });

    _audioService.playbackStateStream.listen((state) {
      if (mounted) {
        setState(() {
          _playbackState = state;
        });

        // Handle album art rotation
        if (state == PlaybackState.playing) {
          _albumArtController.repeat();
        } else {
          _albumArtController.stop();
        }
      }
    });

    _audioService.positionStream.listen((position) {
      if (mounted && !_isDragging) {
        setState(() {
          _position = position;
        });
      }
    });

    _audioService.durationStream.listen((duration) {
      if (mounted) {
        setState(() {
          _duration = duration;
        });
      }
    });

    // Initialize with current state
    setState(() {
      _currentSong = _audioService.currentSong;
      _playbackState = _audioService.isPlaying ? PlaybackState.playing : PlaybackState.paused;
      _position = _audioService.position;
      _duration = _audioService.duration;
      _isFavorite = _currentSong != null ? _storageService.isFavorite(_currentSong!) : false;
    });
  }

  @override
  void dispose() {
    _albumArtController.dispose();
    _slideController.dispose();
    super.dispose();
  }

  Future<void> _togglePlayPause() async {
    if (_playbackState == PlaybackState.playing) {
      await _audioService.pause();
    } else {
      await _audioService.play();
    }
  }

  Future<void> _toggleFavorite() async {
    if (_currentSong == null) return;

    try {
      if (_isFavorite) {
        await _storageService.removeFromFavorites(_currentSong!);
      } else {
        await _storageService.addToFavorites(_currentSong!);
      }
      setState(() {
        _isFavorite = !_isFavorite;
      });

      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(_isFavorite ? 'Added to favorites' : 'Removed from favorites'),
          duration: const Duration(seconds: 1),
        ),
      );
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Error updating favorites: $e')),
      );
    }
  }

  void _onSeek(double value) {
    final position = Duration(milliseconds: (value * _duration.inMilliseconds).round());
    _audioService.seek(position);
  }

  String _formatDuration(Duration duration) {
    String twoDigits(int n) => n.toString().padLeft(2, "0");
    String twoDigitMinutes = twoDigits(duration.inMinutes.remainder(60));
    String twoDigitSeconds = twoDigits(duration.inSeconds.remainder(60));
    return "${twoDigits(duration.inHours)}:$twoDigitMinutes:$twoDigitSeconds";
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.transparent,
      body: Container(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
            colors: [
              const Color(0xFF0D1421).withOpacity(0.9),
              const Color(0xFF1A1A2E).withOpacity(0.9),
              const Color(0xFF16213E).withOpacity(0.9),
            ],
          ),
        ),
        child: SlideTransition(
          position: _slideAnimation,
          child: SafeArea(
            child: Column(
              children: [
                _buildHeader(),
                Expanded(
                  child: _currentSong != null
                      ? _buildPlayerContent()
                      : _buildNoSongPlaying(),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildHeader() {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          IconButton(
            icon: const Icon(Icons.keyboard_arrow_down, color: Colors.white, size: 28),
            onPressed: () => Navigator.of(context).pop(),
          ),
          Column(
            children: [
              Text(
                'NOW PLAYING',
                style: TextStyle(
                  fontSize: 12,
                  fontWeight: FontWeight.w600,
                  color: Colors.white.withOpacity(0.7),
                  letterSpacing: 1.2,
                ),
              ),
              if (_currentSong != null) ...[
                const SizedBox(height: 4),
                Text(
                  'From "${_currentSong!.album}"',
                  style: TextStyle(
                    fontSize: 11,
                    color: Colors.white.withOpacity(0.5),
                  ),
                ),
              ],
            ],
          ),
          PopupMenuButton<String>(
            icon: const Icon(Icons.more_vert, color: Colors.white),
            color: const Color(0xFF1A1A2E),
            onSelected: (value) {
              switch (value) {
                case 'queue':
                  // Show queue
                  break;
                case 'share':
                  // Share song
                  break;
                case 'info':
                  // Show song info
                  break;
              }
            },
            itemBuilder: (context) => [
              const PopupMenuItem(
                value: 'queue',
                child: Row(
                  children: [
                    Icon(Icons.queue_music, color: Colors.white, size: 20),
                    SizedBox(width: 12),
                    Text('View Queue', style: TextStyle(color: Colors.white)),
                  ],
                ),
              ),
              const PopupMenuItem(
                value: 'share',
                child: Row(
                  children: [
                    Icon(Icons.share, color: Colors.white, size: 20),
                    SizedBox(width: 12),
                    Text('Share', style: TextStyle(color: Colors.white)),
                  ],
                ),
              ),
              const PopupMenuItem(
                value: 'info',
                child: Row(
                  children: [
                    Icon(Icons.info, color: Colors.white, size: 20),
                    SizedBox(width: 12),
                    Text('Song Info', style: TextStyle(color: Colors.white)),
                  ],
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildPlayerContent() {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 24.0),
      child: Column(
        children: [
          const SizedBox(height: 20),
          _buildAlbumArt(),
          const SizedBox(height: 40),
          _buildSongInfo(),
          const SizedBox(height: 30),
          _buildProgressBar(),
          const SizedBox(height: 40),
          _buildPlaybackControls(),
          const SizedBox(height: 30),
          _buildBottomActions(),
          const SizedBox(height: 20),
        ],
      ),
    );
  }

  Widget _buildAlbumArt() {
    return Container(
      width: 300,
      height: 300,
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        boxShadow: [
          BoxShadow(
            color: const Color(0xFFE94560).withOpacity(0.3),
            blurRadius: 40,
            spreadRadius: 10,
          ),
        ],
      ),
      child: AnimatedBuilder(
        animation: _albumArtAnimation,
        builder: (context, child) {
          return Transform.rotate(
            angle: _albumArtAnimation.value,
            child: Container(
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                border: Border.all(
                  color: Colors.white.withOpacity(0.1),
                  width: 8,
                ),
              ),
              child: ClipOval(
                child: CachedNetworkImage(
                  imageUrl: _currentSong?.albumArt ?? '',
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
                      Icons.music_note,
                      color: Colors.white,
                      size: 120,
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
                      Icons.music_note,
                      color: Colors.white,
                      size: 120,
                    ),
                  ),
                ),
              ),
            ),
          );
        },
      ),
    );
  }

  Widget _buildSongInfo() {
    return Column(
      children: [
        Text(
          _currentSong?.title ?? 'Unknown Title',
          style: const TextStyle(
            fontSize: 24,
            fontWeight: FontWeight.bold,
            color: Colors.white,
          ),
          textAlign: TextAlign.center,
          maxLines: 2,
          overflow: TextOverflow.ellipsis,
        ),
        const SizedBox(height: 8),
        Text(
          _currentSong?.artist ?? 'Unknown Artist',
          style: TextStyle(
            fontSize: 16,
            color: Colors.white.withOpacity(0.7),
          ),
          textAlign: TextAlign.center,
        ),
      ],
    );
  }

  Widget _buildProgressBar() {
    return Column(
      children: [
        SliderTheme(
          data: SliderTheme.of(context).copyWith(
            activeTrackColor: const Color(0xFFE94560),
            inactiveTrackColor: Colors.white.withOpacity(0.2),
            thumbColor: const Color(0xFFE94560),
            thumbShape: const RoundSliderThumbShape(enabledThumbRadius: 6),
            overlayShape: const RoundSliderOverlayShape(overlayRadius: 16),
            trackHeight: 4,
          ),
          child: Slider(
            value: _duration.inMilliseconds > 0
                ? (_position.inMilliseconds / _duration.inMilliseconds).clamp(0.0, 1.0)
                : 0.0,
            onChanged: (value) {
              setState(() {
                _isDragging = true;
              });
            },
            onChangeEnd: (value) {
              setState(() {
                _isDragging = false;
              });
              _onSeek(value);
            },
          ),
        ),
        const SizedBox(height: 8),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 24.0),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                _formatDuration(_position),
                style: TextStyle(
                  fontSize: 12,
                  color: Colors.white.withOpacity(0.7),
                ),
              ),
              Text(
                _formatDuration(_duration),
                style: TextStyle(
                  fontSize: 12,
                  color: Colors.white.withOpacity(0.7),
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildPlaybackControls() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        IconButton(
          onPressed: _audioService.toggleShuffle,
          icon: Icon(
            Icons.shuffle,
            color: _audioService.isShuffled
                ? const Color(0xFFE94560)
                : Colors.white.withOpacity(0.6),
            size: 28,
          ),
        ),
        IconButton(
          onPressed: _audioService.previous,
          icon: const Icon(
            Icons.skip_previous,
            color: Colors.white,
            size: 40,
          ),
        ),
        Container(
          width: 80,
          height: 80,
          decoration: BoxDecoration(
            shape: BoxShape.circle,
            gradient: const LinearGradient(
              colors: [
                Color(0xFFE94560),
                Color(0xFFF16E00),
              ],
            ),
            boxShadow: [
              BoxShadow(
                color: const Color(0xFFE94560).withOpacity(0.3),
                blurRadius: 20,
                spreadRadius: 2,
              ),
            ],
          ),
          child: IconButton(
            onPressed: _togglePlayPause,
            icon: Icon(
              _playbackState == PlaybackState.playing
                  ? Icons.pause
                  : Icons.play_arrow,
              color: Colors.white,
              size: 36,
            ),
          ),
        ),
        IconButton(
          onPressed: _audioService.next,
          icon: const Icon(
            Icons.skip_next,
            color: Colors.white,
            size: 40,
          ),
        ),
        IconButton(
          onPressed: _audioService.toggleRepeatMode,
          icon: Icon(
            _audioService.repeatMode == RepeatMode.one
                ? Icons.repeat_one
                : Icons.repeat,
            color: _audioService.repeatMode != RepeatMode.off
                ? const Color(0xFFE94560)
                : Colors.white.withOpacity(0.6),
            size: 28,
          ),
        ),
      ],
    );
  }

  Widget _buildBottomActions() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: [
        IconButton(
          onPressed: _toggleFavorite,
          icon: Icon(
            _isFavorite ? Icons.favorite : Icons.favorite_outline,
            color: _isFavorite ? const Color(0xFFE94560) : Colors.white.withOpacity(0.6),
            size: 28,
          ),
        ),
        IconButton(
          onPressed: () {
            // Show lyrics or queue
          },
          icon: Icon(
            Icons.queue_music,
            color: Colors.white.withOpacity(0.6),
            size: 28,
          ),
        ),
        IconButton(
          onPressed: () {
            // Show share options
          },
          icon: Icon(
            Icons.share,
            color: Colors.white.withOpacity(0.6),
            size: 28,
          ),
        ),
        IconButton(
          onPressed: () {
            // Download song
          },
          icon: Icon(
            Icons.download,
            color: Colors.white.withOpacity(0.6),
            size: 28,
          ),
        ),
      ],
    );
  }

  Widget _buildNoSongPlaying() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            width: 120,
            height: 120,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              color: Colors.white.withOpacity(0.1),
              border: Border.all(
                color: Colors.white.withOpacity(0.2),
                width: 2,
              ),
            ),
            child: const Icon(
              Icons.music_off,
              color: Colors.white54,
              size: 60,
            ),
          ),
          const SizedBox(height: 24),
          const Text(
            'No Song Playing',
            style: TextStyle(
              fontSize: 24,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'Choose a song from your library\nto start playing',
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
}
