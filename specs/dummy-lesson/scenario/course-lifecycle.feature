@scenario
Feature: Instructor creates and publishes a course

  Scenario: Full course lifecycle
    Given POST Register {"Email": "inst@test.com", "Password": "Pass1234!", "Name": "Instructor"} → user
    And POST Login {"Email": "inst@test.com", "Password": "Pass1234!"} → token
    When POST CreateCourse {"Title": "Go 101", "Category": "dev", "Level": "beginner", "Price": 10000} → course
    And POST CreateLesson {"CourseID": course.ID, "Title": "Intro", "VideoURL": "https://example.com/v1", "SortOrder": 1} → lesson
    And PUT PublishCourse {"CourseID": course.ID}
    Then GET ListCourses → courses
    And response.courses contains course.ID
    And status == 200
