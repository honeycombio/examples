package io.honeycomb.examples.javaotlp;

import io.honeycomb.libhoney.*;
import io.honeycomb.libhoney.responses.*;

import java.net.*;

import static io.honeycomb.libhoney.LibHoney.*;

public class Honey {
    private static HoneyClient honeyClientInstance;
    final static ResponseObserver responseObserver = new ResponseObserver() {
        @Override
        public void onServerAccepted(ServerAccepted serverAccepted) {
            System.out.println("onServerAccepted: " + serverAccepted.getMessage());
        }

        @Override
        public void onServerRejected(ServerRejected serverRejected) {
            System.out.println("onServerRejected: " + serverRejected.getMessage());
        }

        @Override
        public void onClientRejected(ClientRejected clientRejected) {
            System.out.println("onClientRejected: " + clientRejected.getMessage());
        }

        @Override
        public void onUnknown(Unknown unknown) {
            System.out.println("onUnknown: " + unknown.getMessage());
        }
    };

    public static HoneyClient getHoneyClient() {
        if (honeyClientInstance == null) {
            final String apiKey = System.getenv("HONEYCOMB_API_KEY");
            final String endpoint = System.getenv("HONEYCOMB_API_ENDPOINT");
            final String dataset = System.getenv("HONEYCOMB_DATASET");
            honeyClientInstance = new HoneyClient(options()
                .setWriteKey(apiKey)
                .setApiHost(URI.create(endpoint))
                .setDataset(dataset)
                .build());
            honeyClientInstance.addResponseObserver(responseObserver);
        }
        return honeyClientInstance;
    }
}
