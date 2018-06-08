package io.honeycomb.libhoney.example.webapp;

import io.honeycomb.libhoney.example.webapp.instrumentation.HoneycombContext;
import io.honeycomb.libhoney.example.webapp.persistence.Todo;
import io.honeycomb.libhoney.example.webapp.persistence.TodoRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.util.List;
import java.util.function.Supplier;

/**
 * Wraps the Todos repository and times calls the database calls. Metrics on the call time are submitted to
 * Honeycomb.
 *
 * This is a simple way of producing the metrics; in a more complex application, it would be preferable to use
 * AOP. See https://github.com/spring-projects/spring-data-examples/tree/master/jpa/interceptors/src/main/java/example/springdata/jpa/interceptors
 * for example.
 */
@Component
public class TodosService {
    @Autowired
    private TodoRepository todoRepository;
    @Autowired
    private HoneycombContext honeycombContext;

    public List<Todo> readTodos() {
        return timeCall(() -> todoRepository.findAll(), "timers.db.select_all_todos");
    }

    public void deleteTodo(final Long id) {
        timeCall(() -> todoRepository.delete(id), "timers.db.delete_todo");
    }

    public void updateTodo(final Long id, final Todo update) {
        timeCall(() -> {
            final Todo todo = todoRepository.getOne(id);
            todo.setCompleted(update.getCompleted());
            todo.setDescription(update.getDescription());
            todo.setDue(update.getDue());
            todoRepository.save(todo);
        }, "timers.db.update_todo");
    }

    public void createTodo(final Todo todo) {
        timeCall(() -> todoRepository.save(todo), "timers.db.insert_todo");
    }

    private  <T> T timeCall(final Supplier<T> call, final String callName) {
        final long startTime = System.currentTimeMillis();
        try {
            return call.get();
        } finally {
            final long endTime = System.currentTimeMillis();
            honeycombContext.getEvent().addField(callName, endTime - startTime);
        }
    }

    private void timeCall(final Runnable call, final String callName) {
        timeCall(() -> { call.run(); return null; }, callName);
    }
}
