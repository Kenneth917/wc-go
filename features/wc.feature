Feature: Count lines, words and bytes
  As a user piping text through the shell
  I want wc to report counts like the system wc
  So that I can swap it into scripts

  Background:
    Given a file "test.txt" with content:
      """
      the quick brown fox
      jumps over
      the lazy dog
      """

  Scenario: Count lines with -l
    When I run "wc -l test.txt"
    Then the output should be:
      """
      3 test.txt
      """

  Scenario: Count words with -w
    When I run "wc -w test.txt"
    Then the output should be:
      """
      9 test.txt
      """

  Scenario: Count bytes with -c
    When I run "wc -c test.txt"
    Then the output should be:
      """
      44 test.txt
      """

  Scenario: Default output is lines, words, bytes
    When I run "wc test.txt"
    Then the output should be:
      """
      3 9 44 test.txt
      """

  Scenario: Read from standard input when no file is given
    When I run "wc -l" with stdin:
      """
      one
      two
      """
    Then the output should be:
      """
      2
      """

  Scenario: Missing file is an error
    When I run "wc -l nope.txt"
    Then the exit code should be 1
