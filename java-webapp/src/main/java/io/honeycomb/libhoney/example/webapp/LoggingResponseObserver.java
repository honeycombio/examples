package io.honeycomb.libhoney.example.webapp;

import io.honeycomb.libhoney.ResponseObserver;
import io.honeycomb.libhoney.responses.ClientRejected;
import io.honeycomb.libhoney.responses.ServerAccepted;
import io.honeycomb.libhoney.responses.ServerRejected;
import io.honeycomb.libhoney.responses.Unknown;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class LoggingResponseObserver implements ResponseObserver {
    private static final Logger LOG = LoggerFactory.getLogger(LoggingResponseObserver.class);

    @Override
    public void onServerAccepted(ServerAccepted serverAccepted) {
        LOG.info("Server accepted: {}", serverAccepted);
    }

    @Override
    public void onServerRejected(ServerRejected serverRejected) {
        LOG.error("Server rejected: {}", serverRejected);
    }

    @Override
    public void onClientRejected(ClientRejected clientRejected) {
        LOG.info("Client rejected: {}", clientRejected);
    }

    @Override
    public void onUnknown(Unknown unknown) {
        LOG.error("Unknown: {}", unknown);
    }
}
