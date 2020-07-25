mkdir results

CALL test_json.cmd > results\test_json.txt

CALL test_login.cmd > results\test_login.txt

CALL test_jwt.cmd > results\test_jwt.txt

CALL test_flight.cmd > results\test_flight.txt

CALL test_flights.cmd > results\test_flights.txt

CALL test_flights_sorted.cmd > results\test_flights_sorted.txt
