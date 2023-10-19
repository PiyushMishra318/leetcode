import is_palindrome from "../problems/9_is_palindrome";

test("is palindrome", function () {
  expect(is_palindrome(121)).toEqual(true);
  expect(is_palindrome(-121)).toEqual(false);
  expect(is_palindrome(10)).toEqual(false);
});
