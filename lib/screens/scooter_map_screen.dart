import 'package:flutter/material.dart';
import 'package:google_maps_flutter/google_maps_flutter.dart';
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
  GoogleMapController? _mapController;
  Set<Marker> _markers = {};
  List<Scooter> _scooters = [];
  bool _isLoading = true;
  String? _error;

  // Baku city center.
  static const _bakuCenter = LatLng(40.4093, 49.8671);

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
        latitude: _bakuCenter.latitude,
        longitude: _bakuCenter.longitude,
        radiusKm: 5.0,
      );

      setState(() {
        _scooters = scooters;
        _markers = _buildMarkers(scooters);
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = e.toString();
        _isLoading = false;
      });
    }
  }

  Set<Marker> _buildMarkers(List<Scooter> scooters) {
    return scooters.map((scooter) {
      return Marker(
        markerId: MarkerId(scooter.id),
        position: LatLng(scooter.latitude, scooter.longitude),
        icon: _markerIconForStatus(scooter.status),
        onTap: () => _onScooterTapped(scooter),
      );
    }).toSet();
  }

  BitmapDescriptor _markerIconForStatus(ScooterStatus status) {
    switch (status) {
      case ScooterStatus.available:
        return BitmapDescriptor.defaultMarkerWithHue(BitmapDescriptor.hueGreen);
      case ScooterStatus.inUse:
        return BitmapDescriptor.defaultMarkerWithHue(BitmapDescriptor.hueRed);
      case ScooterStatus.maintenance:
        return BitmapDescriptor.defaultMarkerWithHue(
          BitmapDescriptor.hueOrange,
        );
      case ScooterStatus.lowBattery:
        return BitmapDescriptor.defaultMarkerWithHue(
          BitmapDescriptor.hueYellow,
        );
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

  // BUG: No optimistic UI update — user sees stale state until next refresh.
  // BUG: No error handling if unlock fails after bottom sheet closes.
  Future<void> _unlockScooter(Scooter scooter) async {
    Navigator.of(context).pop();

    try {
      await _scooterService.unlockScooter(scooter.id);
      _loadScooters(); // Full reload instead of targeted update.
    } catch (e) {
      // BUG: Error is silently swallowed if context is no longer mounted.
      if (mounted) {
        ScaffoldMessenger.of(
          context,
        ).showSnackBar(SnackBar(content: Text('Unlock failed: $e')));
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Jet Sharing'),
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
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text('Error: $_error', style: const TextStyle(color: Colors.red)),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: _loadScooters,
              child: const Text('Retry'),
            ),
          ],
        ),
      );
    }

    return GoogleMap(
      initialCameraPosition: const CameraPosition(
        target: _bakuCenter,
        zoom: 14.0,
      ),
      markers: _markers,
      myLocationEnabled: true,
      myLocationButtonEnabled: true,
      onMapCreated: (controller) => _mapController = controller,
    );
  }

  @override
  void dispose() {
    _mapController?.dispose();
    super.dispose();
  }
}
