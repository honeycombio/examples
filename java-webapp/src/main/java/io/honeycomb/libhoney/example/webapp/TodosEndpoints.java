package io.honeycomb.libhoney.example.webapp;

import io.honeycomb.libhoney.example.webapp.persistence.Todo;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.util.MimeTypeUtils;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping(value = "todos")
public class TodosEndpoints {
    @Autowired
    private TodosService todosService;

    @GetMapping(produces = MimeTypeUtils.APPLICATION_JSON_VALUE)
    public List<Todo> todos() {
        return todosService.readTodos();
    }

    @PutMapping(value = "/{id}", consumes = MimeTypeUtils.APPLICATION_JSON_VALUE)
    public void put(final @PathVariable Long id, @RequestBody final Todo todoUpdate) {
        todosService.updateTodo(id, todoUpdate);
    }

    @DeleteMapping(value = "/{id}")
    public void delete(final @PathVariable Long id) {
        todosService.deleteTodo(id);
    }

    @PostMapping(consumes = MimeTypeUtils.APPLICATION_JSON_VALUE)
    public void post(@RequestBody final Todo todoUpdate) {
        todosService.createTodo(todoUpdate);
    }
}
