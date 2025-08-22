import 'package:flutter/material.dart';
import '../../models/song.dart';
import '../../services/music_api_service.dart';
import '../../services/audio_service.dart';
import '../../widgets/song_item.dart';
import '../../utils/app_colors.dart';

class GenreScreen extends StatefulWidget {
  const GenreScreen({super.key});

  @override
  State<GenreScreen> createState() => _GenreScreenState();
}

class _GenreScreenState extends State<GenreScreen>
    with TickerProviderStateMixin {
  final MusicApiService _musicApiService = MusicApiService();
  final AudioService _audioService = AudioService();

  List<String> genres = [];
  String selectedGenre = 'All';
  List<Song> songs = [];
  List<Song> filteredSongs = [];
  bool isLoading = true;
  String sortBy = 'title'; // title, artist, duration
  bool sortAscending = true;

  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _loadGenres();
    _loadSongs();
  }

  Future<void> _loadGenres() async {
    try {
      // Mock genres - in real app, get from API
      genres = [
        'All',
        'Pop',
        'Rock',
        'Hip-Hop',
        'R&B',
        'Country',
        'Electronic',
        'Jazz',
        'Classical',
        'Reggae',
        'Folk',
        'Alternative',
        'Indie',
        'Blues',
        'Punk'
      ];

      _tabController = TabController(length: genres.length, vsync: this);
      _tabController.addListener(_onTabChanged);

      setState(() {});
    } catch (e) {
      debugPrint('Error loading genres: $e');
    }
  }

  void _onTabChanged() {
    if (!_tabController.indexIsChanging) {
      setState(() {
        selectedGenre = genres[_tabController.index];
        _filterSongs();
      });
    }
  }

  Future<void> _loadSongs() async {
    setState(() => isLoading = true);

    try {
      final loadedSongs = await _musicApiService.getTrendingSongs();
      
      setState(() {
        songs = loadedSongs;
        _filterSongs();
        isLoading = false;
      });
    } catch (e) {
      debugPrint('Error loading songs: $e');
      setState(() => isLoading = false);
    }
  }

  void _filterSongs() {
    if (selectedGenre == 'All') {
      filteredSongs = List.from(songs);
    } else {
      filteredSongs = songs
          .where((song) => song.genre?.toLowerCase() == selectedGenre.toLowerCase())
          .toList();
    }
    _sortSongs();
  }

  void _sortSongs() {
    filteredSongs.sort((a, b) {
      int comparison;
      switch (sortBy) {
        case 'artist':
          comparison = a.artist.toLowerCase().compareTo(b.artist.toLowerCase());
          break;
        case 'duration':
          comparison = a.duration.inMilliseconds.compareTo(b.duration.inMilliseconds);
          break;
        case 'title':
        default:
          comparison = a.title.toLowerCase().compareTo(b.title.toLowerCase());
          break;
      }
      return sortAscending ? comparison : -comparison;
    });
  }

  void _showSortDialog() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Sort by'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            RadioListTile<String>(
              title: const Text('Title'),
              value: 'title',
              groupValue: sortBy,
              onChanged: (value) => _updateSort(value!),
            ),
            RadioListTile<String>(
              title: const Text('Artist'),
              value: 'artist',
              groupValue: sortBy,
              onChanged: (value) => _updateSort(value!),
            ),
            RadioListTile<String>(
              title: const Text('Duration'),
              value: 'duration',
              groupValue: sortBy,
              onChanged: (value) => _updateSort(value!),
            ),
            const Divider(),
            SwitchListTile(
              title: const Text('Ascending'),
              value: sortAscending,
              onChanged: (value) {
                setState(() {
                  sortAscending = value;
                  _sortSongs();
                });
                Navigator.pop(context);
              },
            ),
          ],
        ),
      ),
    );
  }

  void _updateSort(String newSortBy) {
    setState(() {
      sortBy = newSortBy;
      _sortSongs();
    });
    Navigator.pop(context);
  }

  Future<void> _playAll() async {
    if (filteredSongs.isNotEmpty) {
      await _audioService.playFromPlaylist(filteredSongs, 0);
    }
  }

  Future<void> _shuffleAll() async {
    if (filteredSongs.isNotEmpty) {
      _audioService.toggleShuffle();
      await _audioService.playFromPlaylist(filteredSongs, 0);
    }
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        backgroundColor: AppColors.primary,
        elevation: 0,
        title: const Text(
          'Browse by Genre',
          style: TextStyle(
            fontWeight: FontWeight.bold,
            color: Colors.white,
          ),
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.sort, color: Colors.white),
            onPressed: _showSortDialog,
          ),
          IconButton(
            icon: const Icon(Icons.search, color: Colors.white),
            onPressed: () {
              Navigator.pushNamed(context, '/search');
            },
          ),
        ],
        bottom: genres.isNotEmpty
            ? TabBar(
                controller: _tabController,
                isScrollable: true,
                indicatorColor: AppColors.accent,
                labelColor: Colors.white,
                unselectedLabelColor: Colors.white70,
                tabs: genres.map((genre) => Tab(text: genre)).toList(),
              )
            : null,
      ),
      body: Column(
        children: [
          // Stats and actions bar
          Container(
            padding: const EdgeInsets.all(16),
            color: AppColors.surface,
            child: Row(
              children: [
                Expanded(
                  child: Text(
                    '${filteredSongs.length} songs in $selectedGenre',
                    style: const TextStyle(
                      color: AppColors.textSecondary,
                      fontSize: 14,
                    ),
                  ),
                ),
                if (filteredSongs.isNotEmpty) ...[
                  IconButton(
                    icon: const Icon(Icons.play_arrow, color: AppColors.primary),
                    onPressed: _playAll,
                    tooltip: 'Play All',
                  ),
                  IconButton(
                    icon: const Icon(Icons.shuffle, color: AppColors.primary),
                    onPressed: _shuffleAll,
                    tooltip: 'Shuffle All',
                  ),
                ],
              ],
            ),
          ),

          // Songs list
          Expanded(
            child: isLoading
                ? const Center(
                    child: CircularProgressIndicator(
                      valueColor: AlwaysStoppedAnimation<Color>(AppColors.primary),
                    ),
                  )
                : filteredSongs.isEmpty
                    ? Center(
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Icon(
                              Icons.music_note_outlined,
                              size: 64,
                              color: Colors.grey[600],
                            ),
                            const SizedBox(height: 16),
                            Text(
                              'No songs found in $selectedGenre',
                              style: TextStyle(
                                color: Colors.grey[600],
                                fontSize: 16,
                              ),
                            ),
                            const SizedBox(height: 8),
                            TextButton(
                              onPressed: _loadSongs,
                              child: const Text(
                                'Refresh',
                                style: TextStyle(color: AppColors.primary),
                              ),
                            ),
                          ],
                        ),
                      )
                    : RefreshIndicator(
                        onRefresh: _loadSongs,
                        color: AppColors.primary,
                        child: ListView.builder(
                          itemCount: filteredSongs.length,
                          itemBuilder: (context, index) {
                            final song = filteredSongs[index];
                            return SongItem(
                              song: song,
                              onTap: () async {
                                await _audioService.playFromPlaylist(
                                  filteredSongs,
                                  index,
                                );
                              },
                            );
                          },
                        ),
                      ),
          ),
        ],
      ),
    );
  }
}
