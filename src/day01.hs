{-# LANGUAGE LambdaCase #-}

import Control.Monad (unless)
import Text.Read (readMaybe)
import Data.Maybe (catMaybes)

testDataP1 = "1abc2\npqr3stu8vwx\na1b2c3d4e5f\ntreb7uchet"
testDataP2 = "two1nine\neightwothree\nabcone2threexyz\nxtwone3four\n4nineeightseven2\nzoneight234\n7pqrstsixteen"

isDigit c = c `elem` ['1'..'9']

cleanString :: String -> Maybe Int
cleanString str =
  readMaybe $ head nums : last nums : []
  where
    nums = filter isDigit str


parseDigit :: String -> String
parseDigit = \case
  'z' : 'e' : 'r' : 'o' : xs -> '0' : xs
  'o' : 'r' : 'e' : 'z' : xs -> '0' : xs
  'o' : 'n' : 'e' : xs -> '1':xs
  'e' : 'n' : 'o' : xs -> '1':xs
  't' : 'w' : 'o' : xs -> '2':xs
  'o' : 'w' : 't' : xs -> '2':xs
  't' : 'h' : 'r' : 'e' : 'e' : xs -> '3':xs
  'e' : 'e' : 'r' : 'h' : 't' : xs -> '3':xs
  'f' : 'o' : 'u' : 'r' : xs -> '4':xs
  'r' : 'u' : 'o' : 'f' : xs -> '4':xs
  'f' : 'i' : 'v' : 'e' : xs -> '5':xs
  'e' : 'v' : 'i' : 'f' : xs -> '5':xs
  's' : 'i' : 'x' : xs -> '6':xs
  'x' : 'i' : 's' : xs -> '6':xs
  's' : 'e' : 'v' : 'e' : 'n' : xs -> '7':xs
  'n' : 'e' : 'v' : 'e' : 's' : xs -> '7':xs
  'e' : 'i' : 'g' : 'h' : 't' : xs -> '8':xs
  't' : 'h' : 'g' : 'i' : 'e' : xs -> '8':xs
  'n' : 'i' : 'n' : 'e' : xs -> '9':xs
  'e' : 'n' : 'i' : 'n' : xs -> '9':xs
  x : xs | isDigit x -> x:xs
  x:xs -> parseDigit xs
  _ -> ""


solve1st :: [String] -> Int
solve1st = sum . catMaybes . map cleanString

solve2nd :: [String] -> Int
solve2nd =
  solve1st . map parseDigit_ 
  where
    parseDigit_ = parseDigit . reverse . parseDigit . reverse 

main = do 
  unless ((solve1st $ lines testDataP1) == 142 ) (error "test fail for part one")
  unless ((solve2nd $ lines testDataP2) == 281 ) (error "test fail for part two")

  contents <- readFile "dat/day01.txt"
  let input = lines contents

  putStrLn $ show $ solve1st input
  putStrLn $ show $ solve2nd input
