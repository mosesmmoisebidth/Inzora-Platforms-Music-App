import 'dart:async';
import 'dart:math';
import 'package:just_audio/just_audio.dart';
import 'package:just_audio_background/just_audio_background.dart';
import 'package:audio_session/audio_session.dart';
import '../models/song.dart';

enum PlaybackState {
  stopped,
  playing,
  paused,
  buffering,
  completed,
  error,
}

enum RepeatMode {
  off,
  one,
  all,
}

class AudioService {
  static final AudioService _instance = AudioService._internal();
  factory AudioService() => _instance;
  AudioService._internal();

  final AudioPlayer _audioPlayer = AudioPlayer();
  final StreamController<PlaybackState> _playbackStateController = StreamController<PlaybackState>.broadcast();
  final StreamController<Song?> _currentSongController = StreamController<Song?>.broadcast();
  final StreamController<Duration> _positionController = StreamController<Duration>.broadcast();
  final StreamController<Duration> _durationController = StreamController<Duration>.broadcast();

  List<Song> _playlist = [];
  int _currentIndex = 0;
  bool _isShuffled = false;
  RepeatMode _repeatMode = RepeatMode.off;
  List<int> _shuffledIndices = [];

  // Getters
  Stream<PlaybackState> get playbackStateStream => _playbackStateController.stream;
  Stream<Song?> get currentSongStream => _currentSongController.stream;
  Stream<Duration> get positionStream => _positionController.stream;
  Stream<Duration> get durationStream => _durationController.stream;
  
  Song? get currentSong => _playlist.isEmpty ? null : _playlist[_currentIndex];
  bool get isPlaying => _audioPlayer.playing;
  bool get isShuffled => _isShuffled;
  RepeatMode get repeatMode => _repeatMode;
  List<Song> get playlist => _playlist;
  int get currentIndex => _currentIndex;
  Duration get position => _audioPlayer.position;
  Duration get duration => _audioPlayer.duration ?? Duration.zero;

  Future<void> initialize() async {
    try {
      final session = await AudioSession.instance;
      await session.configure(const AudioSessionConfiguration.music());

      _audioPlayer.playerStateStream.listen((playerState) {
        final processingState = playerState.processingState;
        final playing = playerState.playing;

        PlaybackState state;
        if (processingState == ProcessingState.loading || processingState == ProcessingState.buffering) {
          state = PlaybackState.buffering;
        } else if (!playing) {
          state = PlaybackState.paused;
        } else if (processingState == ProcessingState.completed) {
          state = PlaybackState.completed;
          _onSongCompleted();
        } else {
          state = PlaybackState.playing;
        }

        _playbackStateController.add(state);
      });

      _audioPlayer.positionStream.listen((position) {
        _positionController.add(position);
      });

      _audioPlayer.durationStream.listen((duration) {
        if (duration != null) {
          _durationController.add(duration);
        }
      });
    } catch (e) {
      print('Error initializing audio service: $e');
    }
  }

  Future<void> setPlaylist(List<Song> songs, {int startIndex = 0}) async {
    try {
      _playlist = songs;
      _currentIndex = startIndex;
      _shuffledIndices = List.generate(songs.length, (index) => index);
      
      if (_isShuffled) {
        _shufflePlaylist();
      }

      await _loadCurrentSong();
    } catch (e) {
      _playbackStateController.add(PlaybackState.error);
      print('Error setting playlist: $e');
    }
  }

  Future<void> _loadCurrentSong() async {
    if (_playlist.isEmpty) return;

    try {
      final song = currentSong;
      if (song == null) return;

      _currentSongController.add(song);

      // Create MediaItem for background playback
      final mediaItem = MediaItem(
        id: song.id,
        album: song.album,
        title: song.title,
        artist: song.artist,
        duration: song.duration,
        artUri: Uri.tryParse(song.albumArt),
      );

      // Set audio source with metadata
      await _audioPlayer.setAudioSource(
        AudioSource.uri(
          Uri.parse(song.audioUrl),
          tag: mediaItem,
        ),
      );
    } catch (e) {
      _playbackStateController.add(PlaybackState.error);
      print('Error loading song: $e');
    }
  }

  Future<void> play() async {
    try {
      await _audioPlayer.play();
    } catch (e) {
      _playbackStateController.add(PlaybackState.error);
      print('Error playing: $e');
    }
  }

  Future<void> pause() async {
    try {
      await _audioPlayer.pause();
    } catch (e) {
      print('Error pausing: $e');
    }
  }

  Future<void> stop() async {
    try {
      await _audioPlayer.stop();
      _playbackStateController.add(PlaybackState.stopped);
    } catch (e) {
      print('Error stopping: $e');
    }
  }

  Future<void> seek(Duration position) async {
    try {
      await _audioPlayer.seek(position);
    } catch (e) {
      print('Error seeking: $e');
    }
  }

  Future<void> next() async {
    if (_playlist.isEmpty) return;

    try {
      if (_isShuffled) {
        final currentShuffledIndex = _shuffledIndices.indexOf(_currentIndex);
        if (currentShuffledIndex < _shuffledIndices.length - 1) {
          _currentIndex = _shuffledIndices[currentShuffledIndex + 1];
        } else if (_repeatMode == RepeatMode.all) {
          _currentIndex = _shuffledIndices[0];
        } else {
          return;
        }
      } else {
        if (_currentIndex < _playlist.length - 1) {
          _currentIndex++;
        } else if (_repeatMode == RepeatMode.all) {
          _currentIndex = 0;
        } else {
          return;
        }
      }

      await _loadCurrentSong();
      if (isPlaying) {
        await play();
      }
    } catch (e) {
      print('Error going to next song: $e');
    }
  }

  Future<void> previous() async {
    if (_playlist.isEmpty) return;

    try {
      if (_isShuffled) {
        final currentShuffledIndex = _shuffledIndices.indexOf(_currentIndex);
        if (currentShuffledIndex > 0) {
          _currentIndex = _shuffledIndices[currentShuffledIndex - 1];
        } else if (_repeatMode == RepeatMode.all) {
          _currentIndex = _shuffledIndices[_shuffledIndices.length - 1];
        } else {
          return;
        }
      } else {
        if (_currentIndex > 0) {
          _currentIndex--;
        } else if (_repeatMode == RepeatMode.all) {
          _currentIndex = _playlist.length - 1;
        } else {
          return;
        }
      }

      await _loadCurrentSong();
      if (isPlaying) {
        await play();
      }
    } catch (e) {
      print('Error going to previous song: $e');
    }
  }

  void toggleShuffle() {
    _isShuffled = !_isShuffled;
    if (_isShuffled) {
      _shufflePlaylist();
    }
  }

  void _shufflePlaylist() {
    _shuffledIndices = List.generate(_playlist.length, (index) => index);
    _shuffledIndices.shuffle(Random());
    
    // Make sure current song stays as current
    final currentSongIndex = _shuffledIndices.indexOf(_currentIndex);
    if (currentSongIndex != 0) {
      _shuffledIndices.removeAt(currentSongIndex);
      _shuffledIndices.insert(0, _currentIndex);
    }
  }

  void toggleRepeatMode() {
    switch (_repeatMode) {
      case RepeatMode.off:
        _repeatMode = RepeatMode.all;
        break;
      case RepeatMode.all:
        _repeatMode = RepeatMode.one;
        break;
      case RepeatMode.one:
        _repeatMode = RepeatMode.off;
        break;
    }
  }

  Future<void> _onSongCompleted() async {
    switch (_repeatMode) {
      case RepeatMode.one:
        await seek(Duration.zero);
        await play();
        break;
      case RepeatMode.all:
        await next();
        break;
      case RepeatMode.off:
        if (_currentIndex < _playlist.length - 1) {
          await next();
        }
        break;
    }
  }

  Future<void> playFromPlaylist(List<Song> songs, int index) async {
    await setPlaylist(songs, startIndex: index);
    await play();
  }

  Future<void> playFromSong(Song song) async {
    await setPlaylist([song], startIndex: 0);
    await play();
  }

  Future<void> addToQueue(Song song) async {
    _playlist.add(song);
    if (_isShuffled) {
      _shuffledIndices.add(_playlist.length - 1);
    }
  }

  Future<void> removeFromQueue(int index) async {
    if (index < 0 || index >= _playlist.length) return;
    
    _playlist.removeAt(index);
    
    if (_isShuffled) {
      _shuffledIndices.remove(index);
      // Update indices that are greater than the removed index
      for (int i = 0; i < _shuffledIndices.length; i++) {
        if (_shuffledIndices[i] > index) {
          _shuffledIndices[i]--;
        }
      }
    }
    
    if (_currentIndex >= index) {
      _currentIndex = (_currentIndex - 1).clamp(0, _playlist.length - 1);
    }
  }

  Future<void> setVolume(double volume) async {
    await _audioPlayer.setVolume(volume.clamp(0.0, 1.0));
  }

  Future<void> setSpeed(double speed) async {
    await _audioPlayer.setSpeed(speed.clamp(0.5, 2.0));
  }

  void dispose() {
    _audioPlayer.dispose();
    _playbackStateController.close();
    _currentSongController.close();
    _positionController.close();
    _durationController.close();
  }
}
