import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

import 'package:jet_sharing/main.dart';

void main() {
  testWidgets('App renders Jet Sharing title', (WidgetTester tester) async {
    await tester.pumpWidget(const JetSharingApp());

    expect(find.text('Jet Sharing'), findsOneWidget);
  });
}
