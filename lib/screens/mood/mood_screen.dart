import 'package:flutter/material.dart';
import '../../models/song.dart';
import '../../services/music_api_service.dart';
import '../../services/audio_service.dart';
import '../../widgets/song_item.dart';
import '../../utils/app_colors.dart';

class MoodScreen extends StatefulWidget {
  const MoodScreen({super.key});

  @override
  State<MoodScreen> createState() => _MoodScreenState();
}

class _MoodScreenState extends State<MoodScreen> {
  final MusicApiService _musicApiService = MusicApiService();
  final AudioService _audioService = AudioService();

  List<MoodCategory> moodCategories = [];
  MoodCategory? selectedMood;
  List<Song> songs = [];
  List<Song> filteredSongs = [];
  bool isLoading = true;
  String sortBy = 'title';
  bool sortAscending = true;

  @override
  void initState() {
    super.initState();
    _loadMoodCategories();
    _loadSongs();
  }

  void _loadMoodCategories() {
    // Mock mood categories - in real app, get from API
    moodCategories = [
      MoodCategory(
        id: 'happy',
        name: 'Happy',
        description: 'Uplifting and energetic songs',
        icon: Icons.sentiment_very_satisfied,
        color: Colors.yellow,
        gradient: [Colors.yellow.shade300, Colors.orange.shade400],
      ),
      MoodCategory(
        id: 'chill',
        name: 'Chill',
        description: 'Relaxing and laid-back vibes',
        icon: Icons.self_improvement,
        color: Colors.blue,
        gradient: [Colors.blue.shade300, Colors.teal.shade400],
      ),
      MoodCategory(
        id: 'energetic',
        name: 'Energetic',
        description: 'High energy workout songs',
        icon: Icons.flash_on,
        color: Colors.red,
        gradient: [Colors.red.shade400, Colors.pink.shade400],
      ),
      MoodCategory(
        id: 'romantic',
        name: 'Romantic',
        description: 'Love songs and ballads',
        icon: Icons.favorite,
        color: Colors.pink,
        gradient: [Colors.pink.shade300, Colors.purple.shade400],
      ),
      MoodCategory(
        id: 'melancholy',
        name: 'Melancholy',
        description: 'Emotional and introspective',
        icon: Icons.sentiment_neutral,
        color: Colors.grey,
        gradient: [Colors.grey.shade400, Colors.blueGrey.shade400],
      ),
      MoodCategory(
        id: 'party',
        name: 'Party',
        description: 'Dance and party anthems',
        icon: Icons.celebration,
        color: Colors.purple,
        gradient: [Colors.purple.shade400, Colors.indigo.shade400],
      ),
      MoodCategory(
        id: 'focus',
        name: 'Focus',
        description: 'Instrumental and concentration',
        icon: Icons.psychology,
        color: Colors.green,
        gradient: [Colors.green.shade400, Colors.teal.shade400],
      ),
      MoodCategory(
        id: 'nostalgic',
        name: 'Nostalgic',
        description: 'Classic hits and memories',
        icon: Icons.history,
        color: Colors.amber,
        gradient: [Colors.amber.shade400, Colors.orange.shade400],
      ),
    ];
    setState(() {});
  }

  Future<void> _loadSongs() async {
    setState(() => isLoading = true);

    try {
      final loadedSongs = await _musicApiService.getTrendingSongs();
      
      setState(() {
        songs = loadedSongs;
        if (selectedMood != null) {
          _filterSongs();
        }
        isLoading = false;
      });
    } catch (e) {
      debugPrint('Error loading songs: $e');
      setState(() => isLoading = false);
    }
  }

  void _selectMood(MoodCategory mood) {
    setState(() {
      selectedMood = mood;
      _filterSongs();
    });
  }

  void _filterSongs() {
    if (selectedMood == null) {
      filteredSongs = [];
      return;
    }

    // Mock filtering based on mood - in real app, this would be based on song metadata
    switch (selectedMood!.id) {
      case 'happy':
        filteredSongs = songs.where((song) => 
          song.title.toLowerCase().contains('happy') ||
          song.title.toLowerCase().contains('joy') ||
          song.genre == 'Pop'
        ).toList();
        break;
      case 'chill':
        filteredSongs = songs.where((song) => 
          song.title.toLowerCase().contains('chill') ||
          song.title.toLowerCase().contains('relax') ||
          song.genre == 'R&B'
        ).toList();
        break;
      case 'energetic':
        filteredSongs = songs.where((song) => 
          song.title.toLowerCase().contains('energy') ||
          song.title.toLowerCase().contains('power') ||
          song.genre == 'Electronic'
        ).toList();
        break;
      case 'romantic':
        filteredSongs = songs.where((song) => 
          song.title.toLowerCase().contains('love') ||
          song.title.toLowerCase().contains('heart') ||
          song.title.toLowerCase().contains('kiss')
        ).toList();
        break;
      case 'melancholy':
        filteredSongs = songs.where((song) => 
          song.title.toLowerCase().contains('blue') ||
          song.title.toLowerCase().contains('sad') ||
          song.genre == 'Blues'
        ).toList();
        break;
      case 'party':
        filteredSongs = songs.where((song) => 
          song.title.toLowerCase().contains('dance') ||
          song.title.toLowerCase().contains('party') ||
          song.genre == 'Hip-Hop'
        ).toList();
        break;
      case 'focus':
        filteredSongs = songs.where((song) => 
          song.genre == 'Classical' ||
          song.genre == 'Jazz'
        ).toList();
        break;
      case 'nostalgic':
        filteredSongs = songs.where((song) => 
          song.title.toLowerCase().contains('memory') ||
          song.title.toLowerCase().contains('yesterday') ||
          song.genre == 'Rock'
        ).toList();
        break;
      default:
        filteredSongs = songs.take(10).toList();
    }

    // Add fallback songs if filtered list is empty
    if (filteredSongs.isEmpty) {
      filteredSongs = songs.take(10).toList();
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

  Widget _buildMoodGrid() {
    return GridView.builder(
      padding: const EdgeInsets.all(16),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        mainAxisSpacing: 16,
        crossAxisSpacing: 16,
        childAspectRatio: 1.2,
      ),
      itemCount: moodCategories.length,
      itemBuilder: (context, index) {
        final mood = moodCategories[index];
        return GestureDetector(
          onTap: () => _selectMood(mood),
          child: Container(
            decoration: BoxDecoration(
              gradient: LinearGradient(
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
                colors: mood.gradient,
              ),
              borderRadius: BorderRadius.circular(16),
              boxShadow: [
                BoxShadow(
                  color: mood.color.withOpacity(0.3),
                  blurRadius: 8,
                  offset: const Offset(0, 4),
                ),
              ],
            ),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Icon(
                  mood.icon,
                  size: 48,
                  color: Colors.white,
                ),
                const SizedBox(height: 12),
                Text(
                  mood.name,
                  style: const TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                  ),
                ),
                const SizedBox(height: 4),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 8),
                  child: Text(
                    mood.description,
                    textAlign: TextAlign.center,
                    style: const TextStyle(
                      fontSize: 12,
                      color: Colors.white70,
                    ),
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
              ],
            ),
          ),
        );
      },
    );
  }

  Widget _buildSongsList() {
    return Column(
      children: [
        // Header with back button and actions
        Container(
          padding: const EdgeInsets.all(16),
          color: AppColors.surface,
          child: Row(
            children: [
              IconButton(
                icon: const Icon(Icons.arrow_back, color: AppColors.primary),
                onPressed: () {
                  setState(() {
                    selectedMood = null;
                    filteredSongs = [];
                  });
                },
              ),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      selectedMood!.name,
                      style: const TextStyle(
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                        color: AppColors.textPrimary,
                      ),
                    ),
                    Text(
                      '${filteredSongs.length} songs • ${selectedMood!.description}',
                      style: const TextStyle(
                        color: AppColors.textSecondary,
                        fontSize: 14,
                      ),
                    ),
                  ],
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
          child: filteredSongs.isEmpty
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Icon(
                        selectedMood!.icon,
                        size: 64,
                        color: selectedMood!.color.withOpacity(0.5),
                      ),
                      const SizedBox(height: 16),
                      Text(
                        'No ${selectedMood!.name.toLowerCase()} songs found',
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
              : ListView.builder(
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
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: selectedMood == null
          ? AppBar(
              backgroundColor: AppColors.primary,
              elevation: 0,
              title: const Text(
                'Browse by Mood',
                style: TextStyle(
                  fontWeight: FontWeight.bold,
                  color: Colors.white,
                ),
              ),
              actions: [
                IconButton(
                  icon: const Icon(Icons.search, color: Colors.white),
                  onPressed: () {
                    Navigator.pushNamed(context, '/search');
                  },
                ),
              ],
            )
          : null,
      body: isLoading
          ? const Center(
              child: CircularProgressIndicator(
                valueColor: AlwaysStoppedAnimation<Color>(AppColors.primary),
              ),
            )
          : selectedMood == null
              ? _buildMoodGrid()
              : _buildSongsList(),
    );
  }
}

class MoodCategory {
  final String id;
  final String name;
  final String description;
  final IconData icon;
  final Color color;
  final List<Color> gradient;

  MoodCategory({
    required this.id,
    required this.name,
    required this.description,
    required this.icon,
    required this.color,
    required this.gradient,
  });
}
