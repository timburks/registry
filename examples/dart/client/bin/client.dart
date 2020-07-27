// Copyright 2020 Google LLC. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Dart test client for the Registry API.
import 'dart:async';
import 'dart:io';

import 'package:args/args.dart';
import 'package:grpc/grpc.dart';
import 'package:client/generated/google/cloud/apigee/registry/v1alpha1/registry_service.pb.dart';
import 'package:client/generated/google/cloud/apigee/registry/v1alpha1/registry_service.pbgrpc.dart';

Future<Null> main(List<String> args) async {
  // Create a gRPC channel to a local Registry server.
  final channel = new ClientChannel(
    'localhost',
    port: 8080,
    options: const ChannelOptions(
      credentials: const ChannelCredentials.insecure(),
    ),
  );
  final channelCompleter = Completer<void>();
  ProcessSignal.sigint.watch().listen((_) async {
    print("sigint begin");
    await channel.terminate();
    channelCompleter.complete();
    print("sigint end");
  });
  // Use the channel to create a Registry client.
  final stub = new RegistryClient(channel);
  // Make a sample API request.
  try {
    final request = ListProductsRequest();
    request.parent = "projects/google";
    while (true) {
      final response = await stub.listProducts(request);
      response.products.forEach((product) => print(product.name));
      if (response.products.length == 0) {
        break;
      }
      request.pageToken = response.nextPageToken;
    }
  } catch (e) {
    print('Caught error: $e');
  }
  await channel.shutdown();
  exit(0);
}