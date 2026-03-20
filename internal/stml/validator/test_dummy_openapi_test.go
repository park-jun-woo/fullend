//ff:func feature=stml-validate type=test-helper control=sequence
//ff:what 테스트용 OpenAPI YAML 상수
package validator

const dummyOpenAPI = `
openapi: "3.0.3"
info:
  title: Test API
  version: "1.0.0"
paths:
  /reservations:
    post:
      operationId: CreateReservation
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                RoomID:
                  type: integer
                StartAt:
                  type: string
                EndAt:
                  type: string
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema:
                type: object
                properties:
                  reservation:
                    type: object
  /me/reservations:
    get:
      operationId: ListMyReservations
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema:
                type: object
                properties:
                  reservations:
                    type: array
                    items:
                      type: object
  /reservations/{ReservationID}:
    get:
      operationId: GetReservation
      parameters:
        - name: ReservationID
          in: path
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema:
                type: object
                properties:
                  reservation:
                    type: object
  /login:
    post:
      operationId: Login
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                Email:
                  type: string
                Password:
                  type: string
      responses:
        "200":
          description: ok
`

const infraOpenAPI = `
openapi: "3.0.3"
info:
  title: Test API
  version: "1.0.0"
paths:
  /items:
    get:
      operationId: ListItems
      x-pagination:
        style: offset
        defaultLimit: 20
        maxLimit: 100
      x-sort:
        allowed: [name, created_at]
        default: name
        direction: asc
      x-filter:
        allowed: [status, category]
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema:
                type: object
                properties:
                  items:
                    type: array
                    items:
                      type: object
                  total:
                    type: integer
  /simple:
    get:
      operationId: ListSimple
      responses:
        "200":
          description: ok
          content:
            application/json:
              schema:
                type: object
                properties:
                  items:
                    type: array
                    items:
                      type: object
`
