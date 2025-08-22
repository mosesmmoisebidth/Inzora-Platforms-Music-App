import 'package:flutter/material.dart';

class AppColors {
  // Primary colors
  static const Color primary = Color(0xFF6C63FF);
  static const Color primaryDark = Color(0xFF5A52E5);
  static const Color accent = Color(0xFFFF6B9D);
  static const Color secondary = Color(0xFF3F51B5);

  // Background colors
  static const Color background = Color(0xFF0F0F23);
  static const Color surface = Color(0xFF1A1A2E);
  static const Color surfaceVariant = Color(0xFF16213E);
  
  // Card and container colors
  static const Color cardBackground = Color(0xFF1E1E2E);
  static const Color containerBackground = Color(0xFF252545);

  // Text colors
  static const Color textPrimary = Color(0xFFFFFFFF);
  static const Color textSecondary = Color(0xFFB0B0C3);
  static const Color textTertiary = Color(0xFF8B8B9F);
  static const Color textDisabled = Color(0xFF6B6B7D);

  // Status colors
  static const Color success = Color(0xFF4CAF50);
  static const Color error = Color(0xFFFF5252);
  static const Color warning = Color(0xFFFF9800);
  static const Color info = Color(0xFF2196F3);

  // Button colors
  static const Color buttonPrimary = Color(0xFF6C63FF);
  static const Color buttonSecondary = Color(0xFF3F51B5);
  static const Color buttonDisabled = Color(0xFF4A4A5C);

  // Border colors
  static const Color border = Color(0xFF2A2A3E);
  static const Color borderLight = Color(0xFF3A3A54);
  static const Color borderDark = Color(0xFF1A1A2E);

  // Icon colors
  static const Color iconPrimary = Color(0xFFFFFFFF);
  static const Color iconSecondary = Color(0xFFB0B0C3);
  static const Color iconDisabled = Color(0xFF6B6B7D);
  static const Color iconAccent = Color(0xFF6C63FF);

  // Gradient colors
  static const List<Color> primaryGradient = [
    Color(0xFF6C63FF),
    Color(0xFF3F51B5),
  ];

  static const List<Color> accentGradient = [
    Color(0xFFFF6B9D),
    Color(0xFFFF8A80),
  ];

  static const List<Color> backgroundGradient = [
    Color(0xFF0F0F23),
    Color(0xFF1A1A2E),
  ];

  // Overlay colors
  static const Color overlay = Color(0x80000000);
  static const Color overlayLight = Color(0x40000000);
  static const Color overlayDark = Color(0xB3000000);

  // Shimmer colors for loading states
  static const Color shimmerBase = Color(0xFF2A2A3E);
  static const Color shimmerHighlight = Color(0xFF3A3A54);

  // Player colors
  static const Color playerBackground = Color(0xFF1E1E2E);
  static const Color playerControl = Color(0xFF6C63FF);
  static const Color progressTrack = Color(0xFF3A3A54);
  static const Color progressActive = Color(0xFF6C63FF);

  // Social colors
  static const Color facebook = Color(0xFF1877F2);
  static const Color google = Color(0xFFDB4437);
  static const Color spotify = Color(0xFF1DB954);
  static const Color appleMusic = Color(0xFFFC3C44);

  // Transparent colors
  static const Color transparent = Colors.transparent;
  static const Color white = Colors.white;
  static const Color black = Colors.black;

  // Method to get colors based on theme mode
  static Color getTextColor(bool isDarkMode) {
    return isDarkMode ? textPrimary : Colors.black87;
  }

  static Color getBackgroundColor(bool isDarkMode) {
    return isDarkMode ? background : Colors.white;
  }

  static Color getSurfaceColor(bool isDarkMode) {
    return isDarkMode ? surface : Colors.grey[100]!;
  }
}
