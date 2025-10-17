-- Tạo bảng 'posts' để lưu trữ các bài viết.
CREATE TABLE posts (
  -- id: Khóa chính của bài viết.
  id INT AUTO_INCREMENT PRIMARY KEY,

  -- title: Tiêu đề bài viết, không được rỗng.
  title VARCHAR(255) NOT NULL,

  -- content: Nội dung bài viết, có thể dài nên dùng kiểu TEXT.
  content TEXT NULL,

  -- author_id: Khóa ngoại để liên kết với bảng 'users'.
  -- Mỗi bài viết phải thuộc về một user.
  author_id INT NOT NULL,

  -- Dấu thời gian tương tự bảng users.
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  -- Định nghĩa ràng buộc khóa ngoại.
  -- Tên ràng buộc 'fk_posts_author' giúp dễ dàng nhận biết.
  -- Liên kết cột 'author_id' của bảng này tới cột 'id' của bảng 'users'.
  -- ON DELETE CASCADE: Nếu một user bị xóa, tất cả các bài viết của họ cũng sẽ tự động bị xóa.
  CONSTRAINT fk_posts_author
    FOREIGN KEY (author_id)
    REFERENCES users(id)
    ON DELETE CASCADE
) ENGINE=InnoDB;