import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';
import '../models/song.dart';
import '../models/playlist.dart';
import '../models/user.dart';

class StorageService {
  static final StorageService _instance = StorageService._internal();
  factory StorageService() => _instance;
  StorageService._internal();

  late SharedPreferences _prefs;
  bool _initialized = false;

  // Keys for storage
  static const String _userKey = 'current_user';
  static const String _favoriteSongsKey = 'favorite_songs';
  static const String _recentlyPlayedKey = 'recently_played';
  static const String _playlistsKey = 'user_playlists';
  static const String _offlineSongsKey = 'offline_songs';
  static const String _settingsKey = 'app_settings';

  Future<void> initialize() async {
    if (_initialized) return;
    
    _prefs = await SharedPreferences.getInstance();
    _initialized = true;
  }

  void _ensureInitialized() {
    if (!_initialized) {
      throw Exception('StorageService not initialized. Call initialize() first.');
    }
  }

  // User management
  Future<void> saveUser(User user) async {
    _ensureInitialized();
    await _prefs.setString(_userKey, jsonEncode(user.toJson()));
  }

  User? getUser() {
    _ensureInitialized();
    final userJson = _prefs.getString(_userKey);
    if (userJson == null) return null;
    
    try {
      return User.fromJson(jsonDecode(userJson));
    } catch (e) {
      print('Error parsing user data: $e');
      return null;
    }
  }

  Future<void> clearUser() async {
    _ensureInitialized();
    await _prefs.remove(_userKey);
  }

  // Favorite songs
  Future<void> saveFavoriteSongs(List<Song> songs) async {
    _ensureInitialized();
    final songsJson = songs.map((song) => song.toJson()).toList();
    await _prefs.setString(_favoriteSongsKey, jsonEncode(songsJson));
  }

  List<Song> getFavoriteSongs() {
    _ensureInitialized();
    final songsJson = _prefs.getString(_favoriteSongsKey);
    if (songsJson == null) return [];

    try {
      final List<dynamic> songsList = jsonDecode(songsJson);
      return songsList.map((songJson) => Song.fromJson(songJson)).toList();
    } catch (e) {
      print('Error parsing favorite songs: $e');
      return [];
    }
  }

  Future<void> addToFavorites(Song song) async {
    final favorites = getFavoriteSongs();
    if (!favorites.any((s) => s.id == song.id)) {
      favorites.add(song);
      await saveFavoriteSongs(favorites);
    }
  }

  Future<void> removeFromFavorites(Song song) async {
    final favorites = getFavoriteSongs();
    favorites.removeWhere((s) => s.id == song.id);
    await saveFavoriteSongs(favorites);
  }

  bool isFavorite(Song song) {
    final favorites = getFavoriteSongs();
    return favorites.any((s) => s.id == song.id);
  }

  // Recently played songs
  Future<void> addToRecentlyPlayed(Song song) async {
    _ensureInitialized();
    final recentlyPlayed = getRecentlyPlayed();
    
    // Remove if already exists
    recentlyPlayed.removeWhere((s) => s.id == song.id);
    
    // Add to beginning
    recentlyPlayed.insert(0, song);
    
    // Limit to 50 songs
    if (recentlyPlayed.length > 50) {
      recentlyPlayed.removeRange(50, recentlyPlayed.length);
    }
    
    final songsJson = recentlyPlayed.map((song) => song.toJson()).toList();
    await _prefs.setString(_recentlyPlayedKey, jsonEncode(songsJson));
  }

  List<Song> getRecentlyPlayed() {
    _ensureInitialized();
    final songsJson = _prefs.getString(_recentlyPlayedKey);
    if (songsJson == null) return [];

    try {
      final List<dynamic> songsList = jsonDecode(songsJson);
      return songsList.map((songJson) => Song.fromJson(songJson)).toList();
    } catch (e) {
      print('Error parsing recently played songs: $e');
      return [];
    }
  }

  Future<void> clearRecentlyPlayed() async {
    _ensureInitialized();
    await _prefs.remove(_recentlyPlayedKey);
  }

  // User playlists
  Future<void> saveUserPlaylists(List<Playlist> playlists) async {
    _ensureInitialized();
    final playlistsJson = playlists.map((playlist) => playlist.toJson()).toList();
    await _prefs.setString(_playlistsKey, jsonEncode(playlistsJson));
  }

  List<Playlist> getUserPlaylists() {
    _ensureInitialized();
    final playlistsJson = _prefs.getString(_playlistsKey);
    if (playlistsJson == null) return [];

    try {
      final List<dynamic> playlistsList = jsonDecode(playlistsJson);
      return playlistsList.map((playlistJson) => Playlist.fromJson(playlistJson)).toList();
    } catch (e) {
      print('Error parsing user playlists: $e');
      return [];
    }
  }

  Future<void> addUserPlaylist(Playlist playlist) async {
    final playlists = getUserPlaylists();
    final existingIndex = playlists.indexWhere((p) => p.id == playlist.id);
    
    if (existingIndex != -1) {
      playlists[existingIndex] = playlist;
    } else {
      playlists.add(playlist);
    }
    
    await saveUserPlaylists(playlists);
  }

  Future<void> removeUserPlaylist(String playlistId) async {
    final playlists = getUserPlaylists();
    playlists.removeWhere((p) => p.id == playlistId);
    await saveUserPlaylists(playlists);
  }

  // Offline songs
  Future<void> saveOfflineSongs(List<Song> songs) async {
    _ensureInitialized();
    final songsJson = songs.map((song) => song.toJson()).toList();
    await _prefs.setString(_offlineSongsKey, jsonEncode(songsJson));
  }

  List<Song> getOfflineSongs() {
    _ensureInitialized();
    final songsJson = _prefs.getString(_offlineSongsKey);
    if (songsJson == null) return [];

    try {
      final List<dynamic> songsList = jsonDecode(songsJson);
      return songsList.map((songJson) => Song.fromJson(songJson)).toList();
    } catch (e) {
      print('Error parsing offline songs: $e');
      return [];
    }
  }

  Future<void> addToOffline(Song song) async {
    final offlineSongs = getOfflineSongs();
    if (!offlineSongs.any((s) => s.id == song.id)) {
      offlineSongs.add(song.copyWith(isOfflineAvailable: true));
      await saveOfflineSongs(offlineSongs);
    }
  }

  Future<void> removeFromOffline(Song song) async {
    final offlineSongs = getOfflineSongs();
    offlineSongs.removeWhere((s) => s.id == song.id);
    await saveOfflineSongs(offlineSongs);
  }

  bool isOfflineAvailable(Song song) {
    final offlineSongs = getOfflineSongs();
    return offlineSongs.any((s) => s.id == song.id);
  }

  // App settings
  Future<void> saveSetting(String key, dynamic value) async {
    _ensureInitialized();
    final settings = getSettings();
    settings[key] = value;
    await _prefs.setString(_settingsKey, jsonEncode(settings));
  }

  Map<String, dynamic> getSettings() {
    _ensureInitialized();
    final settingsJson = _prefs.getString(_settingsKey);
    if (settingsJson == null) return {};

    try {
      return Map<String, dynamic>.from(jsonDecode(settingsJson));
    } catch (e) {
      print('Error parsing settings: $e');
      return {};
    }
  }

  T? getSetting<T>(String key, [T? defaultValue]) {
    final settings = getSettings();
    return settings[key] ?? defaultValue;
  }

  Future<void> clearAllData() async {
    _ensureInitialized();
    await _prefs.clear();
  }

  // Theme settings
  bool get isDarkMode => getSetting('dark_mode', false) ?? false;
  Future<void> setDarkMode(bool value) => saveSetting('dark_mode', value);

  // Audio settings
  double get volume => getSetting('volume', 1.0) ?? 1.0;
  Future<void> setVolume(double value) => saveSetting('volume', value);

  bool get isAutoplayEnabled => getSetting('autoplay_enabled', true) ?? true;
  Future<void> setAutoplayEnabled(bool value) => saveSetting('autoplay_enabled', value);

  // Language settings
  String get language => getSetting('language', 'en') ?? 'en';
  Future<void> setLanguage(String value) => saveSetting('language', value);

  // First time usage
  bool get isFirstTime => getSetting('is_first_time', true) ?? true;
  Future<void> setFirstTime(bool value) => saveSetting('is_first_time', value);
}
