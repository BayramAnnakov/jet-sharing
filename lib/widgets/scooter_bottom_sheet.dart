import 'package:flutter/material.dart';
import '../models/scooter.dart';

class ScooterBottomSheet extends StatelessWidget {
  final Scooter scooter;
  final VoidCallback onUnlock;

  const ScooterBottomSheet({
    super.key,
    required this.scooter,
    required this.onUnlock,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: const BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      padding: const EdgeInsets.all(24),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildHeader(),
          const SizedBox(height: 16),
          _buildBatteryRow(),
          const SizedBox(height: 12),
          _buildPriceRow(),
          const SizedBox(height: 24),
          _buildUnlockButton(context),
          SizedBox(height: MediaQuery.of(context).padding.bottom),
        ],
      ),
    );
  }

  Widget _buildHeader() {
    return Row(
      children: [
        const Icon(Icons.electric_scooter, size: 32),
        const SizedBox(width: 12),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                scooter.name,
                style: const TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                ),
              ),
              Text(
                scooter.id,
                style: TextStyle(fontSize: 14, color: Colors.grey[600]),
              ),
            ],
          ),
        ),
        _buildStatusChip(),
      ],
    );
  }

  Widget _buildStatusChip() {
    final (label, color) = switch (scooter.status) {
      ScooterStatus.available => ('Available', Colors.green),
      ScooterStatus.inUse => ('In Use', Colors.red),
      ScooterStatus.maintenance => ('Maintenance', Colors.orange),
      ScooterStatus.lowBattery => ('Low Battery', Colors.amber),
    };

    return Chip(
      label: Text(label, style: const TextStyle(color: Colors.white)),
      backgroundColor: color,
    );
  }

  Widget _buildBatteryRow() {
    final batteryColor = scooter.batteryLevel >= 50
        ? Colors.green
        : scooter.batteryLevel >= 20
            ? Colors.orange
            : Colors.red;

    return Row(
      children: [
        Icon(Icons.battery_std, color: batteryColor),
        const SizedBox(width: 8),
        Text('${scooter.batteryLevel}%',
            style: TextStyle(fontSize: 16, color: batteryColor)),
        const Spacer(),
        Text(
          '~${(scooter.batteryLevel * 0.4).round()} min ride',
          style: TextStyle(color: Colors.grey[600]),
        ),
      ],
    );
  }

  Widget _buildPriceRow() {
    return Row(
      children: [
        const Icon(Icons.attach_money, color: Colors.blue),
        const SizedBox(width: 8),
        Text(
          'R\$ ${scooter.pricePerMinute.toStringAsFixed(2)}/min',
          style: const TextStyle(fontSize: 16),
        ),
      ],
    );
  }

  Widget _buildUnlockButton(BuildContext context) {
    return SizedBox(
      width: double.infinity,
      height: 52,
      child: ElevatedButton(
        onPressed: scooter.canUnlock ? onUnlock : null,
        style: ElevatedButton.styleFrom(
          backgroundColor: Colors.blue[700],
          foregroundColor: Colors.white,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
        ),
        child: Text(
          scooter.canUnlock ? 'Unlock Scooter' : 'Not Available',
          style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
        ),
      ),
    );
  }
}
