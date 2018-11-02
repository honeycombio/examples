package io.honeycomb.libhoney.example.webapp.instrumentation;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.IncorrectResultSizeDataAccessException;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.ResponseStatus;

@ControllerAdvice
public class GlobalExceptionHandler {
    private static final Logger LOG = LoggerFactory.getLogger(GlobalExceptionHandler.class);

    @Autowired
    private HoneycombContext honeycombContext;

    @ExceptionHandler(IncorrectResultSizeDataAccessException.class)
    @ResponseStatus(HttpStatus.BAD_REQUEST)
    public void entityNotFound(final IncorrectResultSizeDataAccessException ex) {
        LOG.error("Error", ex);
        honeycombContext.getContextualData().put("errors.message", ex.getMessage());
    }

    @ExceptionHandler(Exception.class)
    @ResponseStatus(HttpStatus.INTERNAL_SERVER_ERROR)
    public void catchAll(final Exception ex) {
        LOG.error("Error", ex);
        honeycombContext.getContextualData().put("errors.message", ex.getMessage());
    }
}
