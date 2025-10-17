-- Bắt đầu lệnh tạo bảng 'users'
CREATE TABLE users (
  -- id: Số nguyên, tự động tăng, và là khóa chính
  id INT AUTO_INCREMENT PRIMARY KEY,

  -- email: Chuỗi, phải là duy nhất và không được để trống (dùng để đăng nhập)
  email VARCHAR(255) UNIQUE NOT NULL,

  -- password: Chuỗi, không được để trống (sẽ lưu mật khẩu đã được hash)
  password VARCHAR(255) NOT NULL,

  -- created_at: Dấu thời gian, tự động đặt là thời gian hiện tại khi tạo mới
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  -- updated_at: Dấu thời gian, tự động cập nhật khi bản ghi được thay đổi
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;
-- ENGINE=InnoDB là standard engine cho MySQL, hỗ trợ transactions.