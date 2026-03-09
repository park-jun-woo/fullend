@scenario
Feature: Student enrolls in a published course

  Background:
    Given POST Register {"Email": "inst@test.com", "Password": "Pass1234!", "Name": "Inst"} → instructor
    And POST Login {"Email": "inst@test.com", "Password": "Pass1234!"} → token
    And POST CreateCourse {"Title": "Go 101", "Category": "dev", "Level": "beginner", "Price": 10000} → course
    And PUT PublishCourse {"CourseID": course.ID}

  Scenario: Successful enrollment
    Given POST Register {"Email": "student@test.com", "Password": "Pass1234!", "Name": "Student"} → student
    And POST Login {"Email": "student@test.com", "Password": "Pass1234!"} → token
    When POST EnrollCourse {"CourseID": course.ID, "PaymentMethod": "card"} → enrollment
    Then status == 200
    And response.enrollment exists
    And response.payment exists
    And GET ListMyEnrollments → myEnrollments
    And response.enrollments contains enrollment.ID
