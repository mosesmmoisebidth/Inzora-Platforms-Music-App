import 'dart:convert';
import 'dart:math';
import 'package:http/http.dart' as http;
import '../models/song.dart';
import '../models/playlist.dart';

class MusicApiService {
  static const String _baseUrl = 'https://jsonplaceholder.typicode.com';
  
  // Mock data for demonstration - In a real app, you'd use services like Spotify, Apple Music, etc.
  final http.Client _client = http.Client();

  // Generate mock songs for demonstration
  List<Song> _generateMockSongs() {
    final List<String> artists = [
      'The Weeknd', 'Billie Eilish', 'Drake', 'Taylor Swift', 'Ed Sheeran',
      'Ariana Grande', 'Post Malone', 'Dua Lipa', 'Justin Bieber', 'Olivia Rodrigo',
      'Harry Styles', 'Bad Bunny', 'Bruno Mars', 'Adele', 'The Chainsmokers'
    ];

    final List<String> songTitles = [
      'Blinding Lights', 'Bad Guy', 'God\'s Plan', 'Anti-Hero', 'Shape of You',
      '7 rings', 'Circles', 'Levitating', 'Peaches', 'Good 4 U',
      'As It Was', 'Me Porto Bonito', 'That\'s What I Like', 'Easy On Me', 'Closer',
      'Starboy', 'Lovely', 'One Dance', 'Shake It Off', 'Perfect',
      'Thank U, Next', 'Rockstar', 'Don\'t Start Now', 'Sorry', 'Drivers License'
    ];

    final List<String> albums = [
      'After Hours', 'When We All Fall Asleep', 'Scorpion', 'Midnights', '÷',
      'Thank U, Next', 'Hollywood\'s Bleeding', 'Future Nostalgia', 'Justice', 'SOUR',
      'Harry\'s House', 'Un Verano Sin Ti', '24K Magic', '30', 'Memories'
    ];

    final List<String> genres = [
      'Pop', 'Hip Hop', 'R&B', 'Electronic', 'Rock', 'Indie', 'Reggaeton', 'Soul'
    ];

    return List.generate(100, (index) {
      final random = Random(index); // Use index as seed for consistent results
      return Song(
        id: 'song_$index',
        title: songTitles[index % songTitles.length],
        artist: artists[index % artists.length],
        album: albums[index % albums.length],
        albumArt: 'https://picsum.photos/300/300?random=$index',
        audioUrl: 'https://www.soundjay.com/misc/sounds-relax.mp3', // Demo URL
        duration: Duration(
          seconds: 120 + random.nextInt(240), // 2-6 minutes
        ),
        genre: genres[random.nextInt(genres.length)],
        playCount: random.nextInt(1000000),
        isFavorite: random.nextBool(),
      );
    });
  }

  List<Playlist> _generateMockPlaylists() {
    final mockSongs = _generateMockSongs();
    final playlistNames = [
      'Today\'s Top Hits', 'Pop Rising', 'Hip Hop Central', 'Chill Vibes',
      'Workout Playlist', 'Indie Mix', 'R&B Classics', 'Electronic Beats',
      'Acoustic Sessions', 'Night Drive', 'Study Focus', 'Party Mix'
    ];

    return List.generate(12, (index) {
      final random = Random(index);
      final playlistSongs = mockSongs
          .skip(index * 8)
          .take(8 + random.nextInt(12))
          .toList();

      return Playlist(
        id: 'playlist_$index',
        name: playlistNames[index],
        description: 'A curated playlist of ${playlistNames[index].toLowerCase()}',
        coverImage: 'https://picsum.photos/300/300?random=${100 + index}',
        songs: playlistSongs,
        createdAt: DateTime.now().subtract(Duration(days: index * 7)),
        updatedAt: DateTime.now().subtract(Duration(days: index)),
        createdBy: 'Music App',
        isPublic: true,
      );
    });
  }

  Future<List<Song>> fetchTrendingSongs({int limit = 20}) async {
    try {
      // Simulate API delay
      await Future.delayed(const Duration(milliseconds: 500));
      
      final mockSongs = _generateMockSongs();
      return mockSongs.take(limit).toList();
    } catch (e) {
      throw Exception('Failed to fetch trending songs: $e');
    }
  }

  // Alias method for backward compatibility
  Future<List<Song>> getTrendingSongs({int limit = 20}) async {
    return fetchTrendingSongs(limit: limit);
  }

  Future<List<Song>> fetchTopCharts({int limit = 50}) async {
    try {
      await Future.delayed(const Duration(milliseconds: 700));
      
      final mockSongs = _generateMockSongs();
      mockSongs.sort((a, b) => b.playCount.compareTo(a.playCount));
      return mockSongs.take(limit).toList();
    } catch (e) {
      throw Exception('Failed to fetch top charts: $e');
    }
  }

  Future<List<Playlist>> fetchFeaturedPlaylists({int limit = 10}) async {
    try {
      await Future.delayed(const Duration(milliseconds: 600));
      
      final mockPlaylists = _generateMockPlaylists();
      return mockPlaylists.take(limit).toList();
    } catch (e) {
      throw Exception('Failed to fetch featured playlists: $e');
    }
  }

  Future<List<Song>> searchSongs(String query, {int limit = 20}) async {
    try {
      await Future.delayed(const Duration(milliseconds: 300));
      
      if (query.isEmpty) return [];
      
      final mockSongs = _generateMockSongs();
      final filteredSongs = mockSongs.where((song) =>
          song.title.toLowerCase().contains(query.toLowerCase()) ||
          song.artist.toLowerCase().contains(query.toLowerCase()) ||
          song.album.toLowerCase().contains(query.toLowerCase())).toList();
      
      return filteredSongs.take(limit).toList();
    } catch (e) {
      throw Exception('Failed to search songs: $e');
    }
  }

  Future<List<Playlist>> searchPlaylists(String query, {int limit = 10}) async {
    try {
      await Future.delayed(const Duration(milliseconds: 300));
      
      if (query.isEmpty) return [];
      
      final mockPlaylists = _generateMockPlaylists();
      final filteredPlaylists = mockPlaylists.where((playlist) =>
          playlist.name.toLowerCase().contains(query.toLowerCase()) ||
          playlist.description.toLowerCase().contains(query.toLowerCase())).toList();
      
      return filteredPlaylists.take(limit).toList();
    } catch (e) {
      throw Exception('Failed to search playlists: $e');
    }
  }

  Future<List<Song>> fetchSongsByGenre(String genre, {int limit = 30}) async {
    try {
      await Future.delayed(const Duration(milliseconds: 400));
      
      final mockSongs = _generateMockSongs();
      final genreSongs = mockSongs.where((song) =>
          song.genre.toLowerCase() == genre.toLowerCase()).toList();
      
      return genreSongs.take(limit).toList();
    } catch (e) {
      throw Exception('Failed to fetch songs by genre: $e');
    }
  }

  Future<List<Song>> fetchRecommendedSongs(List<String> favoriteGenres, {int limit = 20}) async {
    try {
      await Future.delayed(const Duration(milliseconds: 500));
      
      final mockSongs = _generateMockSongs();
      final recommended = <Song>[];
      
      // Add songs from favorite genres
      for (final genre in favoriteGenres) {
        final genreSongs = mockSongs.where((song) =>
            song.genre.toLowerCase() == genre.toLowerCase()).toList();
        recommended.addAll(genreSongs.take(limit ~/ favoriteGenres.length));
      }
      
      // Fill remaining with random songs
      if (recommended.length < limit) {
        final remaining = mockSongs.where((song) => !recommended.contains(song)).toList();
        recommended.addAll(remaining.take(limit - recommended.length));
      }
      
      return recommended.take(limit).toList();
    } catch (e) {
      throw Exception('Failed to fetch recommended songs: $e');
    }
  }

  Future<Playlist> fetchPlaylistDetails(String playlistId) async {
    try {
      await Future.delayed(const Duration(milliseconds: 400));
      
      final mockPlaylists = _generateMockPlaylists();
      return mockPlaylists.firstWhere(
        (playlist) => playlist.id == playlistId,
        orElse: () => throw Exception('Playlist not found'),
      );
    } catch (e) {
      throw Exception('Failed to fetch playlist details: $e');
    }
  }

  void dispose() {
    _client.close();
  }
}
