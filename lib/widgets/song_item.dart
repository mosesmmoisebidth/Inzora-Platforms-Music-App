import 'package:flutter/material.dart';
import '../models/song.dart';
import '../services/audio_service.dart';
import '../utils/app_colors.dart';

class SongItem extends StatefulWidget {
  final Song song;
  final VoidCallback? onTap;
  final bool showDuration;
  final Widget? trailing;

  const SongItem({
    Key? key,
    required this.song,
    this.onTap,
    this.showDuration = true,
    this.trailing,
  }) : super(key: key);

  @override
  State<SongItem> createState() => _SongItemState();
}

class _SongItemState extends State<SongItem> {
  final AudioService _audioService = AudioService();

  @override
  Widget build(BuildContext context) {
    return ListTile(
      contentPadding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
      leading: ClipRRect(
        borderRadius: BorderRadius.circular(8.0),
        child: widget.song.albumArt != null && widget.song.albumArt!.isNotEmpty
            ? Image.network(
                widget.song.albumArt!,
                width: 56,
                height: 56,
                fit: BoxFit.cover,
                errorBuilder: (context, error, stackTrace) {
                  return Container(
                    width: 56,
                    height: 56,
                    color: AppColors.surface,
                    child: const Icon(
                      Icons.music_note,
                      color: AppColors.textSecondary,
                      size: 24,
                    ),
                  );
                },
              )
            : Container(
                width: 56,
                height: 56,
                color: AppColors.surface,
                child: const Icon(
                  Icons.music_note,
                  color: AppColors.textSecondary,
                  size: 24,
                ),
              ),
      ),
      title: Text(
        widget.song.title,
        style: const TextStyle(
          color: AppColors.textPrimary,
          fontSize: 16,
          fontWeight: FontWeight.w500,
        ),
        maxLines: 1,
        overflow: TextOverflow.ellipsis,
      ),
      subtitle: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            widget.song.artist,
            style: const TextStyle(
              color: AppColors.textSecondary,
              fontSize: 14,
            ),
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
          if (widget.song.album != null && widget.song.album!.isNotEmpty)
            Text(
              widget.song.album!,
              style: const TextStyle(
                color: AppColors.textSecondary,
                fontSize: 12,
              ),
              maxLines: 1,
              overflow: TextOverflow.ellipsis,
            ),
        ],
      ),
      trailing: widget.trailing ??
          Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              if (widget.showDuration && widget.song.duration != null)
                Text(
                  _formatDuration(widget.song.duration!),
                  style: const TextStyle(
                    color: AppColors.textSecondary,
                    fontSize: 12,
                  ),
                ),
              const SizedBox(width: 8),
              IconButton(
                onPressed: () {
                  _audioService.playFromSong(widget.song);
                },
                icon: const Icon(
                  Icons.play_arrow,
                  color: AppColors.primary,
                ),
              ),
            ],
          ),
      onTap: widget.onTap ??
          () {
            _audioService.playFromSong(widget.song);
          },
    );
  }

  String _formatDuration(Duration duration) {
    String twoDigits(int n) => n.toString().padLeft(2, "0");
    String twoDigitMinutes = twoDigits(duration.inMinutes.remainder(60));
    String twoDigitSeconds = twoDigits(duration.inSeconds.remainder(60));
    return "${duration.inHours > 0 ? '${duration.inHours}:' : ''}$twoDigitMinutes:$twoDigitSeconds";
  }
}
