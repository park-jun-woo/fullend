@scenario
Feature: Unauthorized access is denied

  Scenario: Enroll without auth
    When POST EnrollCourse {"CourseID": 1, "PaymentMethod": "card"}
    Then status == 401

  Scenario: Update course by non-owner
    Given POST Register {"Email": "other@test.com", "Password": "Pass1234!", "Name": "Other"} → user
    And POST Login {"Email": "other@test.com", "Password": "Pass1234!"} → token
    When PUT UpdateCourse {"CourseID": 1, "Title": "Hacked", "Category": "x", "Level": "x", "Price": 0}
    Then status == 403
