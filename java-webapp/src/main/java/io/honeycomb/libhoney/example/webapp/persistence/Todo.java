package io.honeycomb.libhoney.example.webapp.persistence;

import com.fasterxml.jackson.annotation.JsonFormat;

import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.Id;
import java.util.Date;

@Entity
public class Todo {
    @Id
    @GeneratedValue
    private Long id;

    @Column(nullable = false)
    private String description;

    @Column(nullable = false)
    private Boolean completed;

    @JsonFormat(shape = JsonFormat.Shape.STRING)
    @Column
    private Date due;

    public Todo() {
        //noargs
     }

    public Todo(String description, Boolean completed, Date due) {
        this.description = description;
        this.completed = completed;
        this.due = due;
    }

    public Long getId() {
        return id;
    }

    public String getDescription() {
        return description;
    }

    public Boolean getCompleted() {
        return completed;
    }

    public Date getDue() {
        return due;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public void setCompleted(Boolean completed) {
        this.completed = completed;
    }

    public void setDue(Date due) {
        this.due = due;
    }
}
