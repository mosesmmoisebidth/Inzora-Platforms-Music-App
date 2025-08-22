class User {
  final String id;
  final String email;
  final String displayName;
  final String profileImageUrl;
  final DateTime createdAt;
  final DateTime lastLoginAt;
  final List<String> favoriteGenres;
  final Map<String, dynamic> preferences;
  final bool isPremium;

  User({
    required this.id,
    required this.email,
    required this.displayName,
    this.profileImageUrl = '',
    required this.createdAt,
    required this.lastLoginAt,
    this.favoriteGenres = const [],
    this.preferences = const {},
    this.isPremium = false,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id']?.toString() ?? '',
      email: json['email'] ?? '',
      displayName: json['display_name'] ?? json['name'] ?? '',
      profileImageUrl: json['profile_image_url'] ?? json['photo_url'] ?? '',
      createdAt: json['created_at'] != null
          ? DateTime.tryParse(json['created_at']) ?? DateTime.now()
          : DateTime.now(),
      lastLoginAt: json['last_login_at'] != null
          ? DateTime.tryParse(json['last_login_at']) ?? DateTime.now()
          : DateTime.now(),
      favoriteGenres: List<String>.from(json['favorite_genres'] ?? []),
      preferences: Map<String, dynamic>.from(json['preferences'] ?? {}),
      isPremium: json['is_premium'] ?? false,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'email': email,
      'display_name': displayName,
      'profile_image_url': profileImageUrl,
      'created_at': createdAt.toIso8601String(),
      'last_login_at': lastLoginAt.toIso8601String(),
      'favorite_genres': favoriteGenres,
      'preferences': preferences,
      'is_premium': isPremium,
    };
  }

  User copyWith({
    String? id,
    String? email,
    String? displayName,
    String? profileImageUrl,
    DateTime? createdAt,
    DateTime? lastLoginAt,
    List<String>? favoriteGenres,
    Map<String, dynamic>? preferences,
    bool? isPremium,
  }) {
    return User(
      id: id ?? this.id,
      email: email ?? this.email,
      displayName: displayName ?? this.displayName,
      profileImageUrl: profileImageUrl ?? this.profileImageUrl,
      createdAt: createdAt ?? this.createdAt,
      lastLoginAt: lastLoginAt ?? this.lastLoginAt,
      favoriteGenres: favoriteGenres ?? this.favoriteGenres,
      preferences: preferences ?? this.preferences,
      isPremium: isPremium ?? this.isPremium,
    );
  }

  // Preference getters and setters
  bool get isDarkMode => preferences['dark_mode'] ?? false;
  bool get isNotificationsEnabled => preferences['notifications_enabled'] ?? true;
  String get language => preferences['language'] ?? 'en';
  double get audioQuality => preferences['audio_quality']?.toDouble() ?? 1.0;
  bool get isAutoplayEnabled => preferences['autoplay_enabled'] ?? true;
  bool get isShuffleDefault => preferences['shuffle_default'] ?? false;

  User updatePreference(String key, dynamic value) {
    final newPreferences = Map<String, dynamic>.from(preferences);
    newPreferences[key] = value;
    return copyWith(preferences: newPreferences);
  }

  User addFavoriteGenre(String genre) {
    if (favoriteGenres.contains(genre)) return this;
    return copyWith(favoriteGenres: [...favoriteGenres, genre]);
  }

  User removeFavoriteGenre(String genre) {
    return copyWith(
      favoriteGenres: favoriteGenres.where((g) => g != genre).toList(),
    );
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is User && other.id == id;
  }

  @override
  int get hashCode => id.hashCode;
}
