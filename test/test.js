// a function that reverse a string
function reverseString(str) {
  return str.split('').reverse().join('');
}


// a function that sluggifies a string
function sluggifyString(str) {
    return str
        .toLowerCase()
        .trim()
        .replace(/[^a-z0-9\s-]/g, '')
        .replace(/\s+/g, '-')
        .replace(/-+/g, '-');
}

// Example usage of sluggifyString function
const title = "Hello World! This is a Test.";
const slug = sluggifyString(title);
console.log(slug); // Output: "hello-world-this-is-a-test"

