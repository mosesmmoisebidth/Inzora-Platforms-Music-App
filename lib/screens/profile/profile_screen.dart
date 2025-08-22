import 'package:flutter/material.dart';
import 'package:cached_network_image/cached_network_image.dart';
import '../../models/user.dart';
import '../../services/storage_service.dart';
import '../auth/login_screen.dart';

class ProfileScreen extends StatefulWidget {
  const ProfileScreen({super.key});

  @override
  State<ProfileScreen> createState() => _ProfileScreenState();
}

class _ProfileScreenState extends State<ProfileScreen>
    with AutomaticKeepAliveClientMixin {
  final StorageService _storageService = StorageService();

  User? _currentUser;
  bool _isLoading = true;
  bool _isDarkMode = false;
  bool _isNotificationsEnabled = true;
  double _audioQuality = 1.0;
  bool _isAutoplayEnabled = true;

  @override
  bool get wantKeepAlive => true;

  @override
  void initState() {
    super.initState();
    _loadUserData();
  }

  Future<void> _loadUserData() async {
    setState(() {
      _isLoading = true;
    });

    try {
      final user = _storageService.getUser();
      final isDarkMode = _storageService.isDarkMode;
      final isNotificationsEnabled = _storageService.getSetting('notifications_enabled', true);
      final audioQuality = _storageService.getSetting('audio_quality', 1.0);
      final isAutoplayEnabled = _storageService.isAutoplayEnabled;

      if (mounted) {
        setState(() {
          _currentUser = user;
          _isDarkMode = isDarkMode;
          _isNotificationsEnabled = isNotificationsEnabled ?? true;
          _audioQuality = audioQuality ?? 1.0;
          _isAutoplayEnabled = isAutoplayEnabled;
          _isLoading = false;
        });
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  Future<void> _logout() async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        backgroundColor: const Color(0xFF1A1A2E),
        title: const Text('Logout', style: TextStyle(color: Colors.white)),
        content: const Text(
          'Are you sure you want to logout?',
          style: TextStyle(color: Colors.white70),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.of(context).pop(true),
            child: const Text('Logout'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      try {
        await _storageService.clearUser();
        if (mounted) {
          Navigator.of(context).pushAndRemoveUntil(
            MaterialPageRoute(builder: (context) => const LoginScreen()),
            (route) => false,
          );
        }
      } catch (e) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error logging out: $e')),
        );
      }
    }
  }

  Future<void> _toggleDarkMode(bool value) async {
    setState(() {
      _isDarkMode = value;
    });
    await _storageService.setDarkMode(value);
  }

  Future<void> _toggleNotifications(bool value) async {
    setState(() {
      _isNotificationsEnabled = value;
    });
    await _storageService.saveSetting('notifications_enabled', value);
  }

  Future<void> _updateAudioQuality(double value) async {
    setState(() {
      _audioQuality = value;
    });
    await _storageService.saveSetting('audio_quality', value);
  }

  Future<void> _toggleAutoplay(bool value) async {
    setState(() {
      _isAutoplayEnabled = value;
    });
    await _storageService.setAutoplayEnabled(value);
  }

  @override
  Widget build(BuildContext context) {
    super.build(context);
    
    return Scaffold(
      backgroundColor: Colors.transparent,
      body: SafeArea(
        child: _isLoading
            ? const Center(child: CircularProgressIndicator())
            : SingleChildScrollView(
                padding: const EdgeInsets.all(20.0),
                child: Column(
                  children: [
                    _buildProfileHeader(),
                    const SizedBox(height: 32),
                    _buildSettingsSection(),
                    const SizedBox(height: 24),
                    _buildAccountSection(),
                    const SizedBox(height: 24),
                    _buildAboutSection(),
                    const SizedBox(height: 100),
                  ],
                ),
              ),
      ),
    );
  }

  Widget _buildProfileHeader() {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.05),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(
          color: Colors.white.withOpacity(0.1),
          width: 1,
        ),
      ),
      child: Column(
        children: [
          Container(
            width: 100,
            height: 100,
            decoration: BoxDecoration(
              shape: BoxShape.circle,
              gradient: const LinearGradient(
                colors: [
                  Color(0xFFE94560),
                  Color(0xFFF16E00),
                ],
              ),
              border: Border.all(
                color: Colors.white.withOpacity(0.2),
                width: 3,
              ),
            ),
            child: _currentUser?.profileImageUrl.isNotEmpty == true
                ? ClipOval(
                    child: CachedNetworkImage(
                      imageUrl: _currentUser!.profileImageUrl,
                      fit: BoxFit.cover,
                      placeholder: (context, url) => const CircularProgressIndicator(),
                      errorWidget: (context, url, error) => _buildDefaultAvatar(),
                    ),
                  )
                : _buildDefaultAvatar(),
          ),
          const SizedBox(height: 16),
          Text(
            _currentUser?.displayName ?? 'Music Lover',
            style: const TextStyle(
              fontSize: 24,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            _currentUser?.email ?? 'user@example.com',
            style: TextStyle(
              fontSize: 16,
              color: Colors.white.withOpacity(0.7),
            ),
          ),
          const SizedBox(height: 16),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: [
              _buildStatItem('Songs', '${_storageService.getFavoriteSongs().length}'),
              _buildStatItem('Playlists', '${_storageService.getUserPlaylists().length}'),
              _buildStatItem(
                _currentUser?.isPremium == true ? 'Premium' : 'Free',
                _currentUser?.isPremium == true ? 'Active' : 'Plan',
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildDefaultAvatar() {
    return const Icon(
      Icons.person,
      size: 50,
      color: Colors.white,
    );
  }

  Widget _buildStatItem(String label, String value) {
    return Column(
      children: [
        Text(
          value,
          style: const TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: Color(0xFFE94560),
          ),
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: TextStyle(
            fontSize: 12,
            color: Colors.white.withOpacity(0.6),
          ),
        ),
      ],
    );
  }

  Widget _buildSettingsSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Playback Settings',
          style: TextStyle(
            fontSize: 20,
            fontWeight: FontWeight.bold,
            color: Colors.white,
          ),
        ),
        const SizedBox(height: 16),
        _buildSettingsCard([
          _buildSwitchTile(
            icon: Icons.autorenew,
            title: 'Autoplay',
            subtitle: 'Automatically play similar songs',
            value: _isAutoplayEnabled,
            onChanged: _toggleAutoplay,
          ),
          const Divider(color: Colors.white12),
          _buildSliderTile(
            icon: Icons.high_quality,
            title: 'Audio Quality',
            subtitle: _getQualityText(_audioQuality),
            value: _audioQuality,
            min: 0.5,
            max: 1.0,
            divisions: 2,
            onChanged: _updateAudioQuality,
          ),
        ]),
      ],
    );
  }

  Widget _buildAccountSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'Account & Privacy',
          style: TextStyle(
            fontSize: 20,
            fontWeight: FontWeight.bold,
            color: Colors.white,
          ),
        ),
        const SizedBox(height: 16),
        _buildSettingsCard([
          _buildSwitchTile(
            icon: Icons.notifications,
            title: 'Notifications',
            subtitle: 'Get updates about new music',
            value: _isNotificationsEnabled,
            onChanged: _toggleNotifications,
          ),
          const Divider(color: Colors.white12),
          _buildSwitchTile(
            icon: Icons.dark_mode,
            title: 'Dark Mode',
            subtitle: 'Use dark theme (Coming soon)',
            value: _isDarkMode,
            onChanged: null, // Disabled for now
          ),
          const Divider(color: Colors.white12),
          _buildTile(
            icon: Icons.security,
            title: 'Privacy Settings',
            subtitle: 'Manage your privacy preferences',
            onTap: () {
              // Navigate to privacy settings
            },
          ),
          const Divider(color: Colors.white12),
          _buildTile(
            icon: Icons.download,
            title: 'Download Settings',
            subtitle: 'Manage offline downloads',
            onTap: () {
              // Navigate to download settings
            },
          ),
        ]),
      ],
    );
  }

  Widget _buildAboutSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          'About',
          style: TextStyle(
            fontSize: 20,
            fontWeight: FontWeight.bold,
            color: Colors.white,
          ),
        ),
        const SizedBox(height: 16),
        _buildSettingsCard([
          _buildTile(
            icon: Icons.help,
            title: 'Help & Support',
            subtitle: 'Get help or contact support',
            onTap: () {
              // Navigate to help
            },
          ),
          const Divider(color: Colors.white12),
          _buildTile(
            icon: Icons.info,
            title: 'About MusicFlow',
            subtitle: 'Version 1.0.0',
            onTap: () {
              _showAboutDialog();
            },
          ),
          const Divider(color: Colors.white12),
          _buildTile(
            icon: Icons.star,
            title: 'Rate the App',
            subtitle: 'Help us improve with your feedback',
            onTap: () {
              // Open app store rating
            },
          ),
          const Divider(color: Colors.white12),
          _buildTile(
            icon: Icons.logout,
            title: 'Logout',
            subtitle: 'Sign out of your account',
            onTap: _logout,
            textColor: Colors.red,
          ),
        ]),
      ],
    );
  }

  Widget _buildSettingsCard(List<Widget> children) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(0.05),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: Colors.white.withOpacity(0.1),
          width: 1,
        ),
      ),
      child: Column(children: children),
    );
  }

  Widget _buildTile({
    required IconData icon,
    required String title,
    required String subtitle,
    VoidCallback? onTap,
    Color? textColor,
  }) {
    return ListTile(
      leading: Icon(
        icon,
        color: textColor ?? Colors.white.withOpacity(0.8),
      ),
      title: Text(
        title,
        style: TextStyle(
          fontSize: 16,
          fontWeight: FontWeight.w600,
          color: textColor ?? Colors.white,
        ),
      ),
      subtitle: Text(
        subtitle,
        style: TextStyle(
          fontSize: 14,
          color: (textColor ?? Colors.white).withOpacity(0.7),
        ),
      ),
      trailing: onTap != null
          ? Icon(
              Icons.arrow_forward_ios,
              size: 16,
              color: Colors.white.withOpacity(0.5),
            )
          : null,
      onTap: onTap,
    );
  }

  Widget _buildSwitchTile({
    required IconData icon,
    required String title,
    required String subtitle,
    required bool value,
    ValueChanged<bool>? onChanged,
  }) {
    return ListTile(
      leading: Icon(
        icon,
        color: Colors.white.withOpacity(0.8),
      ),
      title: Text(
        title,
        style: const TextStyle(
          fontSize: 16,
          fontWeight: FontWeight.w600,
          color: Colors.white,
        ),
      ),
      subtitle: Text(
        subtitle,
        style: TextStyle(
          fontSize: 14,
          color: Colors.white.withOpacity(0.7),
        ),
      ),
      trailing: Switch(
        value: value,
        onChanged: onChanged,
        activeColor: const Color(0xFFE94560),
        inactiveThumbColor: Colors.grey,
        inactiveTrackColor: Colors.grey.withOpacity(0.3),
      ),
    );
  }

  Widget _buildSliderTile({
    required IconData icon,
    required String title,
    required String subtitle,
    required double value,
    required double min,
    required double max,
    int? divisions,
    required ValueChanged<double> onChanged,
  }) {
    return Column(
      children: [
        ListTile(
          leading: Icon(
            icon,
            color: Colors.white.withOpacity(0.8),
          ),
          title: Text(
            title,
            style: const TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: Colors.white,
            ),
          ),
          subtitle: Text(
            subtitle,
            style: TextStyle(
              fontSize: 14,
              color: Colors.white.withOpacity(0.7),
            ),
          ),
        ),
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16.0),
          child: Slider(
            value: value,
            min: min,
            max: max,
            divisions: divisions,
            activeColor: const Color(0xFFE94560),
            inactiveColor: Colors.white.withOpacity(0.3),
            onChanged: onChanged,
          ),
        ),
      ],
    );
  }

  String _getQualityText(double quality) {
    if (quality <= 0.5) return 'Normal (96 kbps)';
    if (quality <= 0.75) return 'High (160 kbps)';
    return 'Very High (320 kbps)';
  }

  void _showAboutDialog() {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        backgroundColor: const Color(0xFF1A1A2E),
        title: const Text('About MusicFlow', style: TextStyle(color: Colors.white)),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'MusicFlow - Your Music, Your Way',
              style: TextStyle(
                color: Colors.white,
                fontSize: 16,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 16),
            Text(
              'Version: 1.0.0\n'
              'Build: ${DateTime.now().year}.${DateTime.now().month}\n\n'
              'A beautiful and intuitive music streaming app designed to provide the best listening experience.',
              style: TextStyle(
                color: Colors.white.withOpacity(0.8),
                height: 1.5,
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Close'),
          ),
        ],
      ),
    );
  }
}
