import 'song.dart';

class Playlist {
  final String id;
  final String name;
  final String description;
  final String coverImage;
  final List<Song> songs;
  final DateTime createdAt;
  final DateTime updatedAt;
  final String createdBy;
  final bool isPublic;
  final bool isOfflineAvailable;
  final int totalDuration;

  Playlist({
    required this.id,
    required this.name,
    this.description = '',
    this.coverImage = '',
    required this.songs,
    required this.createdAt,
    required this.updatedAt,
    this.createdBy = '',
    this.isPublic = false,
    this.isOfflineAvailable = false,
    int? totalDuration,
  }) : totalDuration = totalDuration ?? _calculateTotalDuration(songs);

  static int _calculateTotalDuration(List<Song> songs) {
    return songs.fold(0, (total, song) => total + song.duration.inSeconds);
  }

  factory Playlist.fromJson(Map<String, dynamic> json) {
    final songsList = (json['songs'] as List<dynamic>? ?? [])
        .map((songJson) => Song.fromJson(songJson as Map<String, dynamic>))
        .toList();

    return Playlist(
      id: json['id']?.toString() ?? '',
      name: json['name'] ?? 'Unnamed Playlist',
      description: json['description'] ?? '',
      coverImage: json['cover_image'] ?? json['image'] ?? '',
      songs: songsList,
      createdAt: json['created_at'] != null
          ? DateTime.tryParse(json['created_at']) ?? DateTime.now()
          : DateTime.now(),
      updatedAt: json['updated_at'] != null
          ? DateTime.tryParse(json['updated_at']) ?? DateTime.now()
          : DateTime.now(),
      createdBy: json['created_by'] ?? '',
      isPublic: json['is_public'] ?? false,
      isOfflineAvailable: json['is_offline_available'] ?? false,
      totalDuration: json['total_duration'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'description': description,
      'cover_image': coverImage,
      'songs': songs.map((song) => song.toJson()).toList(),
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
      'created_by': createdBy,
      'is_public': isPublic,
      'is_offline_available': isOfflineAvailable,
      'total_duration': totalDuration,
    };
  }

  Playlist copyWith({
    String? id,
    String? name,
    String? description,
    String? coverImage,
    List<Song>? songs,
    DateTime? createdAt,
    DateTime? updatedAt,
    String? createdBy,
    bool? isPublic,
    bool? isOfflineAvailable,
  }) {
    return Playlist(
      id: id ?? this.id,
      name: name ?? this.name,
      description: description ?? this.description,
      coverImage: coverImage ?? this.coverImage,
      songs: songs ?? this.songs,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? DateTime.now(),
      createdBy: createdBy ?? this.createdBy,
      isPublic: isPublic ?? this.isPublic,
      isOfflineAvailable: isOfflineAvailable ?? this.isOfflineAvailable,
    );
  }

  Playlist addSong(Song song) {
    if (songs.contains(song)) return this;
    return copyWith(songs: [...songs, song]);
  }

  Playlist removeSong(Song song) {
    return copyWith(songs: songs.where((s) => s.id != song.id).toList());
  }

  Playlist reorderSongs(int oldIndex, int newIndex) {
    final newSongs = List<Song>.from(songs);
    final song = newSongs.removeAt(oldIndex);
    newSongs.insert(newIndex, song);
    return copyWith(songs: newSongs);
  }

  String get totalDurationText {
    final duration = Duration(seconds: totalDuration);
    final hours = duration.inHours;
    final minutes = duration.inMinutes % 60;
    
    if (hours > 0) {
      return '${hours}h ${minutes}m';
    }
    return '${minutes}m';
  }

  String get songCountText {
    return '${songs.length} song${songs.length != 1 ? 's' : ''}';
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is Playlist && other.id == id;
  }

  @override
  int get hashCode => id.hashCode;
}
