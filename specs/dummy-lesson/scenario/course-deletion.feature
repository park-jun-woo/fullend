@invariant
Feature: Deleted course disappears from listing

  Scenario: Course excluded after deletion
    Given POST Register {"Email": "inst@test.com", "Password": "Pass1234!", "Name": "Inst"} → user
    And POST Login {"Email": "inst@test.com", "Password": "Pass1234!"} → token
    And POST CreateCourse {"Title": "Temp", "Category": "dev", "Level": "beginner", "Price": 0} → course
    When DELETE DeleteCourse {"CourseID": course.ID}
    Then status == 200
    And GET ListCourses → listing
    And response.courses excludes course.ID
