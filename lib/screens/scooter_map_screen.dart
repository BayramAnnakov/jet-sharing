import 'package:flutter/material.dart';
import '../models/scooter.dart';
import '../services/scooter_service.dart';
import '../widgets/scooter_bottom_sheet.dart';

class ScooterMapScreen extends StatefulWidget {
  const ScooterMapScreen({super.key});

  @override
  State<ScooterMapScreen> createState() => _ScooterMapScreenState();
}

class _ScooterMapScreenState extends State<ScooterMapScreen> {
  final ScooterService _scooterService = ScooterService();
  List<Scooter> _scooters = [];
  bool _isLoading = true;
  String? _error;

  @override
  void initState() {
    super.initState();
    _loadScooters();
  }

  Future<void> _loadScooters() async {
    setState(() {
      _isLoading = true;
      _error = null;
    });

    try {
      final scooters = await _scooterService.fetchNearbyScooters(
        latitude: -23.5505,
        longitude: -46.6340,
        radiusKm: 5.0,
      );

      setState(() {
        _scooters = scooters;
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
    }
  }

  void _onScooterTapped(Scooter scooter) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (_) => ScooterBottomSheet(
        scooter: scooter,
        onUnlock: () => _unlockScooter(scooter),
      ),
    );
  }

  Future<void> _unlockScooter(Scooter scooter) async {
    Navigator.of(context).pop();

    try {
      await _scooterService.unlockScooter(scooter.id);
      _loadScooters();
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Unlock failed: $e')),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Jet Sharing'),
        centerTitle: true,
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: _loadScooters,
          ),
        ],
      ),
      body: _buildBody(),
    );
  }

  Widget _buildBody() {
    if (_isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    if (_error != null) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(Icons.cloud_off, size: 64, color: Colors.grey),
              const SizedBox(height: 16),
              Text(
                'Could not load scooters',
                style: Theme.of(context).textTheme.titleMedium,
              ),
              const SizedBox(height: 8),
              Text(
                _error!,
                style: TextStyle(color: Colors.grey[600], fontSize: 13),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 24),
              FilledButton.icon(
                onPressed: _loadScooters,
                icon: const Icon(Icons.refresh),
                label: const Text('Retry'),
              ),
            ],
          ),
        ),
      );
    }

    if (_scooters.isEmpty) {
      return const Center(child: Text('No scooters nearby'));
    }

    return RefreshIndicator(
      onRefresh: _loadScooters,
      child: ListView.separated(
        padding: const EdgeInsets.all(16),
        itemCount: _scooters.length,
        separatorBuilder: (_, __) => const SizedBox(height: 12),
        itemBuilder: (context, index) {
          final scooter = _scooters[index];
          return _ScooterCard(
            scooter: scooter,
            onTap: () => _onScooterTapped(scooter),
          );
        },
      ),
    );
  }
}

class _ScooterCard extends StatelessWidget {
  final Scooter scooter;
  final VoidCallback onTap;

  const _ScooterCard({required this.scooter, required this.onTap});

  @override
  Widget build(BuildContext context) {
    final (statusLabel, statusColor) = switch (scooter.status) {
      ScooterStatus.available => ('Available', Colors.green),
      ScooterStatus.inUse => ('In Use', Colors.red),
      ScooterStatus.maintenance => ('Maintenance', Colors.orange),
      ScooterStatus.lowBattery => ('Low Battery', Colors.amber),
    };

    final batteryColor = scooter.batteryLevel >= 50
        ? Colors.green
        : scooter.batteryLevel >= 20
            ? Colors.orange
            : Colors.red;

    return Card(
      elevation: 2,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(16),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              Container(
                width: 48,
                height: 48,
                decoration: BoxDecoration(
                  color: statusColor.withValues(alpha: 0.15),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Icon(
                  Icons.electric_scooter,
                  color: statusColor,
                  size: 28,
                ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      scooter.name,
                      style: const TextStyle(
                        fontWeight: FontWeight.w600,
                        fontSize: 16,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      scooter.id,
                      style: TextStyle(color: Colors.grey[500], fontSize: 13),
                    ),
                  ],
                ),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 10,
                      vertical: 4,
                    ),
                    decoration: BoxDecoration(
                      color: statusColor.withValues(alpha: 0.15),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Text(
                      statusLabel,
                      style: TextStyle(
                        color: statusColor,
                        fontSize: 12,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                  const SizedBox(height: 6),
                  Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Icon(Icons.battery_std, size: 16, color: batteryColor),
                      const SizedBox(width: 4),
                      Text(
                        '${scooter.batteryLevel}%',
                        style: TextStyle(
                          color: batteryColor,
                          fontWeight: FontWeight.w600,
                          fontSize: 13,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}
