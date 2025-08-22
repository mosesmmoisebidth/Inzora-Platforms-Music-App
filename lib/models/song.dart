class Song {
  final String id;
  final String title;
  final String artist;
  final String album;
  final String albumArt;
  final String audioUrl;
  final Duration duration;
  final String genre;
  final bool isOfflineAvailable;
  final DateTime? releaseDate;
  final int playCount;
  final bool isFavorite;

  Song({
    required this.id,
    required this.title,
    required this.artist,
    required this.album,
    required this.albumArt,
    required this.audioUrl,
    required this.duration,
    this.genre = '',
    this.isOfflineAvailable = false,
    this.releaseDate,
    this.playCount = 0,
    this.isFavorite = false,
  });

  factory Song.fromJson(Map<String, dynamic> json) {
    return Song(
      id: json['id']?.toString() ?? '',
      title: json['title'] ?? json['name'] ?? 'Unknown Title',
      artist: json['artist'] ?? json['artist_name'] ?? 'Unknown Artist',
      album: json['album'] ?? json['album_name'] ?? 'Unknown Album',
      albumArt: json['image'] ?? json['album_art'] ?? json['artwork_url'] ?? '',
      audioUrl: json['audio_url'] ?? json['stream_url'] ?? json['preview_url'] ?? '',
      duration: Duration(
        seconds: int.tryParse(json['duration']?.toString() ?? '0') ?? 0,
      ),
      genre: json['genre'] ?? '',
      isOfflineAvailable: json['is_offline_available'] ?? false,
      releaseDate: json['release_date'] != null 
          ? DateTime.tryParse(json['release_date']) 
          : null,
      playCount: json['play_count'] ?? 0,
      isFavorite: json['is_favorite'] ?? false,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'artist': artist,
      'album': album,
      'album_art': albumArt,
      'audio_url': audioUrl,
      'duration': duration.inSeconds,
      'genre': genre,
      'is_offline_available': isOfflineAvailable,
      'release_date': releaseDate?.toIso8601String(),
      'play_count': playCount,
      'is_favorite': isFavorite,
    };
  }

  Song copyWith({
    String? id,
    String? title,
    String? artist,
    String? album,
    String? albumArt,
    String? audioUrl,
    Duration? duration,
    String? genre,
    bool? isOfflineAvailable,
    DateTime? releaseDate,
    int? playCount,
    bool? isFavorite,
  }) {
    return Song(
      id: id ?? this.id,
      title: title ?? this.title,
      artist: artist ?? this.artist,
      album: album ?? this.album,
      albumArt: albumArt ?? this.albumArt,
      audioUrl: audioUrl ?? this.audioUrl,
      duration: duration ?? this.duration,
      genre: genre ?? this.genre,
      isOfflineAvailable: isOfflineAvailable ?? this.isOfflineAvailable,
      releaseDate: releaseDate ?? this.releaseDate,
      playCount: playCount ?? this.playCount,
      isFavorite: isFavorite ?? this.isFavorite,
    );
  }

  String get durationText {
    final minutes = duration.inMinutes;
    final seconds = duration.inSeconds % 60;
    return '${minutes.toString().padLeft(2, '0')}:${seconds.toString().padLeft(2, '0')}';
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is Song && other.id == id;
  }

  @override
  int get hashCode => id.hashCode;
}
