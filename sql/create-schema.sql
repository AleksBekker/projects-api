-- To be run within an open database

CREATE TABLE Projects (
    project_id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    status ENUM('Not Started', 'In Progress', 'On Hold', 'Completed', 'Cancelled') NOT NULL,
    priority ENUM('Low', 'Medium', 'High') NOT NULL
);

CREATE TABLE Tags (
    tag_id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL
);

CREATE TABLE ProjectTags (
    project_id INT,
    tag_id INT,
    PRIMARY KEY (project_id, tag_id),
    FOREIGN KEY (project_id) REFERENCES Projects(project_id),
    FOREIGN KEY (tag_id) REFERENCES Tags(tag_id)
);

CREATE TABLE ProjectLinks (
    link_id INT PRIMARY KEY AUTO_INCREMENT,
    project_id INT NOT NULL,
    url VARCHAR(2000) NOT NULL,
    display VARCHAR(255) NOT NULL,
    type ENUM('Documentation', 'Source', 'Demo', 'Reference', 'Other') NOT NULL,
    FOREIGN KEY (project_id) REFERENCES Projects(project_id)
);

