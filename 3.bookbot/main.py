def print_report(stats):
  print("--- Begin report of " + stats["book_path"] + " ---")
  print(str(stats["word_count"]) + " words found in the document.\n")
  sorted_letters = sorted(stats["letter_count"].keys())
  for letter in sorted_letters:
    print("The '" + letter + "' character was found "  + str(stats["letter_count"][letter]) + " times")

  print("--- End report ---")

def main(path_to_file):
  book_stats = {}
  book_stats["book_path"] = path_to_file

  with open(path_to_file, 'r') as f:
    content = f.read()
    # print(content)
    # count the number of words in the file
    word_count = len(content.split())
    book_stats["word_count"] = word_count
    # count the letters in the file
    letter_dict = {}
    for line in content.lower():
      for letter in line:
        if letter.isalpha():
          if letter in letter_dict:
            letter_dict[letter] += 1
          else:
            letter_dict[letter] = 1
    
    book_stats["letter_count"] = letter_dict
  # print report
  print_report(book_stats)

if __name__ == "__main__":
  path_to_file = "books/frankenstein.txt"
  main(path_to_file = path_to_file)

