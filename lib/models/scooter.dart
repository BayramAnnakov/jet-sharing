enum ScooterStatus { available, inUse, maintenance, lowBattery }

class Scooter {
  final String id;
  final String name;
  final double latitude;
  final double longitude;
  final int batteryLevel;
  final ScooterStatus status;
  final double pricePerMinute;

  const Scooter({
    required this.id,
    required this.name,
    required this.latitude,
    required this.longitude,
    required this.batteryLevel,
    required this.status,
    required this.pricePerMinute,
  });

  factory Scooter.fromJson(Map<String, dynamic> json) {
    return Scooter(
      id: json['id'] as String,
      name: json['name'] as String,
      latitude: (json['latitude'] as num).toDouble(),
      longitude: (json['longitude'] as num).toDouble(),
      batteryLevel: json['battery_level'] as int,
      status: _parseStatus(json['status'] as String),
      pricePerMinute: (json['price_per_minute'] as num).toDouble(),
    );
  }

  static ScooterStatus _parseStatus(String status) {
    switch (status) {
      case 'available':
        return ScooterStatus.available;
      case 'in_use':
        return ScooterStatus.inUse;
      case 'maintenance':
        return ScooterStatus.maintenance;
      case 'low_battery':
        return ScooterStatus.lowBattery;
      default:
        return ScooterStatus.maintenance;
    }
  }

  bool get canUnlock =>
      status == ScooterStatus.available && batteryLevel >= 10;

  Map<String, dynamic> toJson() => {
    'id': id,
    'name': name,
    'latitude': latitude,
    'longitude': longitude,
    'battery_level': batteryLevel,
    'status': status.name,
    'price_per_minute': pricePerMinute,
  };
}
