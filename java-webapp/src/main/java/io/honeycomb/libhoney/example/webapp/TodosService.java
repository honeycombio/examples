package io.honeycomb.libhoney.example.webapp;

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

    public List<Todo> readTodos() {
        return todoRepository.findAll();
    }

    public void deleteTodo(final Long id) {
        todoRepository.delete(id);
    }

    public void updateTodo(final Long id, final Todo update) {
        final Todo todo = todoRepository.getOne(id);
        todo.setCompleted(update.getCompleted());
        todo.setDescription(update.getDescription());
        todo.setDue(update.getDue());
        todoRepository.save(todo);
    }

    public void createTodo(final Todo todo) {
        todoRepository.save(todo);
    }
}
