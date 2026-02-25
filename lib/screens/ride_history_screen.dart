import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class RideHistoryScreen extends StatefulWidget {
  RideHistoryScreen({super.key});

  @override
  State<RideHistoryScreen> createState() => _RideHistoryScreenState();
}

class _RideHistoryScreenState extends State<RideHistoryScreen> {
  List<dynamic> _rides = [];
  bool _isLoading = true;

  @override
  void initState() {
    super.initState();
    _fetchRideHistory();
  }

  Future<void> _fetchRideHistory() async {
    final response = await http.get(
      Uri.parse('http://localhost:8080/api/rides?user_id=user-1'),
    );

    if (response.statusCode == 200) {
      final List<dynamic> data = jsonDecode(response.body);
      setState(() {
        _rides = data;
        _isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Ride History'),
      ),
      body: _isLoading
          ? Center(child: CircularProgressIndicator())
          : _rides.isEmpty
              ? Center(child: Text('No rides yet'))
              : ListView.builder(
                  itemCount: _rides.length,
                  itemBuilder: (context, index) {
                    final ride = _rides[index];
                    return ListTile(
                      leading: Icon(Icons.electric_scooter),
                      title: Text('Ride ${ride['id']}'),
                      subtitle: Text(
                        'Duration: ${ride['duration']}s  •  Cost: \$${ride['cost']?.toStringAsFixed(2) ?? '0.00'}',
                      ),
                      trailing: _buildStatusChip(ride['status']),
                    );
                  },
                ),
    );
  }

  Widget _buildStatusChip(String? status) {
    final color = switch (status) {
      'completed' => Colors.green,
      'active' => Colors.blue,
      'cancelled' => Colors.red,
      _ => Colors.grey,
    };

    return Chip(
      label: Text(
        status ?? 'unknown',
        style: TextStyle(color: Colors.white, fontSize: 12),
      ),
      backgroundColor: color,
    );
  }
}
