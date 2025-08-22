import 'package:flutter/material.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:flutter_staggered_animations/flutter_staggered_animations.dart';
import '../../models/song.dart';
import '../../models/playlist.dart';
import '../../services/music_api_service.dart';
import '../../services/audio_service.dart';
import '../../services/storage_service.dart';

class TrendingScreen extends StatefulWidget {
  const TrendingScreen({super.key});

  @override
  State<TrendingScreen> createState() => _TrendingScreenState();
}

class _TrendingScreenState extends State<TrendingScreen>
    with SingleTickerProviderStateMixin {
  final MusicApiService _apiService = MusicApiService();
  final AudioService _audioService = AudioService();
  final StorageService _storageService = StorageService();

  late TabController _tabController;
  
  List<Song> _topCharts = [];
  List<Song> _trending = [];
  List<Playlist> _trendingPlaylists = [];
  
  bool _isLoadingCharts = true;
  bool _isLoadingTrending = true;
  bool _isLoadingPlaylists = true;
  
  Song? _currentSong;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    _loadData();
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

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  Future<void> _loadData() async {
    final futures = [
      _loadTopCharts(),
      _loadTrending(),
      _loadTrendingPlaylists(),
    ];

    await Future.wait(futures);
  }

  Future<void> _loadTopCharts() async {
    try {
      final charts = await _apiService.fetchTopCharts(limit: 50);
      if (mounted) {
        setState(() {
          _topCharts = charts;
          _isLoadingCharts = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _isLoadingCharts = false;
        });
      }
    }
  }

  Future<void> _loadTrending() async {
    try {
      final trending = await _apiService.fetchTrendingSongs(limit: 30);
      if (mounted) {
        setState(() {
          _trending = trending;
          _isLoadingTrending = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _isLoadingTrending = false;
        });
      }
    }
  }

  Future<void> _loadTrendingPlaylists() async {
    try {
      final playlists = await _apiService.fetchFeaturedPlaylists(limit: 12);
      if (mounted) {
        setState(() {
          _trendingPlaylists = playlists;
          _isLoadingPlaylists = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _isLoadingPlaylists = false;
        });
      }
    }
  }

  Future<void> _playSong(Song song, List<Song> playlist, int index) async {
    try {
      await _audioService.playFromPlaylist(playlist, index);
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
              _buildTabBar(),
              Expanded(
                child: TabBarView(
                  controller: _tabController,
                  children: [
                    _buildTopChartsTab(),
                    _buildTrendingTab(),
                    _buildTrendingPlaylistsTab(),
                  ],
                ),
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
      child: Row(
        children: [
          IconButton(
            onPressed: () => Navigator.of(context).pop(),
            icon: const Icon(Icons.arrow_back, color: Colors.white),
          ),
          const Expanded(
            child: Text(
              'Trending Music',
              style: TextStyle(
                fontSize: 28,
                fontWeight: FontWeight.bold,
                color: Colors.white,
              ),
            ),
          ),
          IconButton(
            onPressed: () {
              _loadData();
            },
            icon: const Icon(Icons.refresh, color: Colors.white),
          ),
        ],
      ),
    );
  }

  Widget _buildTabBar() {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 20),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
      ),
      child: TabBar(
        controller: _tabController,
        indicator: BoxDecoration(
          color: const Color(0xFFE94560),
          borderRadius: BorderRadius.circular(12),
        ),
        labelColor: Colors.white,
        unselectedLabelColor: Colors.white.withOpacity(0.6),
        labelStyle: const TextStyle(fontWeight: FontWeight.w600),
        dividerColor: Colors.transparent,
        tabs: const [
          Tab(text: 'Top Charts'),
          Tab(text: 'Trending'),
          Tab(text: 'Playlists'),
        ],
      ),
    );
  }

  Widget _buildTopChartsTab() {
    if (_isLoadingCharts) {
      return const Center(child: CircularProgressIndicator());
    }

    if (_topCharts.isEmpty) {
      return _buildEmptyState('No charts available', Icons.trending_up);
    }

    return AnimationLimiter(
      child: ListView.builder(
        padding: const EdgeInsets.all(20),
        itemCount: _topCharts.length,
        itemBuilder: (context, index) {
          return AnimationConfiguration.staggeredList(
            position: index,
            duration: const Duration(milliseconds: 375),
            child: SlideAnimation(
              verticalOffset: 50.0,
              child: FadeInAnimation(
                child: _buildRankedSongTile(_topCharts[index], index + 1, _topCharts),
              ),
            ),
          );
        },
      ),
    );
  }

  Widget _buildTrendingTab() {
    if (_isLoadingTrending) {
      return const Center(child: CircularProgressIndicator());
    }

    if (_trending.isEmpty) {
      return _buildEmptyState('No trending music', Icons.trending_up);
    }

    return AnimationLimiter(
      child: ListView.builder(
        padding: const EdgeInsets.all(20),
        itemCount: _trending.length,
        itemBuilder: (context, index) {
          return AnimationConfiguration.staggeredList(
            position: index,
            duration: const Duration(milliseconds: 375),
            child: SlideAnimation(
              verticalOffset: 50.0,
              child: FadeInAnimation(
                child: _buildSongTile(_trending[index], index, _trending),
              ),
            ),
          );
        },
      ),
    );
  }

  Widget _buildTrendingPlaylistsTab() {
    if (_isLoadingPlaylists) {
      return const Center(child: CircularProgressIndicator());
    }

    if (_trendingPlaylists.isEmpty) {
      return _buildEmptyState('No trending playlists', Icons.playlist_play);
    }

    return AnimationLimiter(
      child: GridView.builder(
        padding: const EdgeInsets.all(20),
        gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
          crossAxisCount: 2,
          childAspectRatio: 0.75,
          crossAxisSpacing: 16,
          mainAxisSpacing: 16,
        ),
        itemCount: _trendingPlaylists.length,
        itemBuilder: (context, index) {
          return AnimationConfiguration.staggeredGrid(
            position: index,
            duration: const Duration(milliseconds: 375),
            columnCount: 2,
            child: ScaleAnimation(
              child: FadeInAnimation(
                child: _buildPlaylistCard(_trendingPlaylists[index]),
              ),
            ),
          );
        },
      ),
    );
  }

  Widget _buildRankedSongTile(Song song, int rank, List<Song> playlist) {
    final isCurrentSong = _currentSong?.id == song.id;
    final isFavorite = _storageService.isFavorite(song);
    final rankColor = rank <= 3 ? const Color(0xFFE94560) : Colors.white.withOpacity(0.7);

    return Container(
      margin: const EdgeInsets.only(bottom: 8),
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
        leading: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 30,
              height: 30,
              decoration: BoxDecoration(
                color: rank <= 3 ? const Color(0xFFE94560) : Colors.transparent,
                shape: BoxShape.circle,
                border: rank > 3 ? Border.all(color: Colors.white.withOpacity(0.3)) : null,
              ),
              child: Center(
                child: Text(
                  rank.toString(),
                  style: TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.bold,
                    color: rank <= 3 ? Colors.white : rankColor,
                  ),
                ),
              ),
            ),
            const SizedBox(width: 12),
            Container(
              width: 50,
              height: 50,
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(8),
              ),
              child: ClipRRect(
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
              '${song.playCount} plays',
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
            IconButton(
              icon: const Icon(Icons.play_arrow, color: Color(0xFFE94560), size: 24),
              onPressed: () => _playSong(song, playlist, rank - 1),
            ),
          ],
        ),
        onTap: () => _playSong(song, playlist, rank - 1),
      ),
    );
  }

  Widget _buildSongTile(Song song, int index, List<Song> playlist) {
    final isCurrentSong = _currentSong?.id == song.id;
    final isFavorite = _storageService.isFavorite(song);

    return Container(
      margin: const EdgeInsets.only(bottom: 8),
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
          child: ClipRRect(
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
        subtitle: Text(
          song.artist,
          style: TextStyle(
            fontSize: 14,
            color: Colors.white.withOpacity(0.7),
          ),
          maxLines: 1,
          overflow: TextOverflow.ellipsis,
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
            IconButton(
              icon: const Icon(Icons.play_arrow, color: Color(0xFFE94560), size: 24),
              onPressed: () => _playSong(song, playlist, index),
            ),
          ],
        ),
        onTap: () => _playSong(song, playlist, index),
      ),
    );
  }

  Widget _buildPlaylistCard(Playlist playlist) {
    return GestureDetector(
      onTap: () {
        // Navigate to playlist detail
      },
      child: Container(
        decoration: BoxDecoration(
          color: Colors.white.withOpacity(0.05),
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
            color: Colors.white.withOpacity(0.1),
            width: 1,
          ),
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Expanded(
              flex: 3,
              child: Container(
                margin: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(8),
                ),
                child: ClipRRect(
                  borderRadius: BorderRadius.circular(8),
                  child: Stack(
                    children: [
                      CachedNetworkImage(
                        imageUrl: playlist.coverImage,
                        fit: BoxFit.cover,
                        width: double.infinity,
                        height: double.infinity,
                        placeholder: (context, url) => Container(
                          decoration: const BoxDecoration(
                            gradient: LinearGradient(
                              colors: [
                                Color(0xFFE94560),
                                Color(0xFFF16E00),
                              ],
                            ),
                          ),
                          child: const Icon(Icons.playlist_play, size: 40, color: Colors.white),
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
                          child: const Icon(Icons.playlist_play, size: 40, color: Colors.white),
                        ),
                      ),
                      Positioned(
                        bottom: 8,
                        right: 8,
                        child: Container(
                          padding: const EdgeInsets.all(6),
                          decoration: BoxDecoration(
                            color: Colors.black.withOpacity(0.7),
                            shape: BoxShape.circle,
                          ),
                          child: const Icon(
                            Icons.play_arrow,
                            color: Colors.white,
                            size: 16,
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
            Expanded(
              flex: 1,
              child: Padding(
                padding: const EdgeInsets.fromLTRB(12, 0, 12, 12),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      playlist.name,
                      style: const TextStyle(
                        fontSize: 14,
                        fontWeight: FontWeight.w600,
                        color: Colors.white,
                      ),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 2),
                    Text(
                      '${playlist.songs.length} songs',
                      style: TextStyle(
                        fontSize: 12,
                        color: Colors.white.withOpacity(0.7),
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildEmptyState(String message, IconData icon) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            icon,
            size: 64,
            color: Colors.white.withOpacity(0.3),
          ),
          const SizedBox(height: 16),
          Text(
            message,
            style: TextStyle(
              fontSize: 18,
              color: Colors.white.withOpacity(0.6),
            ),
          ),
          const SizedBox(height: 32),
          ElevatedButton(
            onPressed: () => _loadData(),
            child: const Text('Retry'),
          ),
        ],
      ),
    );
  }
}
