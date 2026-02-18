import 'package:flutter/material.dart';
import 'screens/scooter_map_screen.dart';

void main() {
  runApp(const JetSharingApp());
}

class JetSharingApp extends StatelessWidget {
  const JetSharingApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Jet Sharing',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        colorSchemeSeed: Colors.blue,
        useMaterial3: true,
      ),
      home: const ScooterMapScreen(),
    );
  }
}
