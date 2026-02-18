import 'dart:convert';
import 'package:http/http.dart' as http;
import '../models/scooter.dart';

class ScooterService {
  // iOS simulator uses localhost to reach the host machine.
  static const _baseUrl = 'http://localhost:8080/api';

  final http.Client _client;

  ScooterService({http.Client? client}) : _client = client ?? http.Client();

  Future<List<Scooter>> fetchNearbyScooters({
    required double latitude,
    required double longitude,
    required double radiusKm,
  }) async {
    final uri = Uri.parse('$_baseUrl/scooters');

    final response = await _client.get(uri, headers: _defaultHeaders());

    if (response.statusCode != 200) {
      throw ScooterServiceException(
        'Failed to fetch scooters: ${response.statusCode}',
        statusCode: response.statusCode,
      );
    }

    final List<dynamic> data = jsonDecode(response.body);
    return data.map((json) => Scooter.fromJson(json)).toList();
  }

  Future<void> unlockScooter(String scooterId) async {
    final uri = Uri.parse('$_baseUrl/scooters/$scooterId/unlock');

    final response = await _client.post(
      uri,
      headers: _defaultHeaders(),
      body: jsonEncode({'action': 'unlock'}),
    );

    if (response.statusCode != 200) {
      throw ScooterServiceException(
        'Failed to unlock scooter: ${response.statusCode}',
        statusCode: response.statusCode,
      );
    }
  }

  // TODO: Read token from flutter_secure_storage instead of hardcoding.
  Map<String, String> _defaultHeaders() => {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer demo-token-placeholder',
  };
}

class ScooterServiceException implements Exception {
  final String message;
  final int? statusCode;

  const ScooterServiceException(this.message, {this.statusCode});

  @override
  String toString() => 'ScooterServiceException: $message (HTTP $statusCode)';
}
