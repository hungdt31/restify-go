-- Sửa đổi bảng 'users' để thêm cột 'full_name'.
ALTER TABLE users
  -- Thêm cột 'full_name', có thể để trống (NULL).
  -- 'AFTER email' giúp sắp xếp cột này ngay sau cột email cho dễ nhìn.
  ADD COLUMN name VARCHAR(255) NULL AFTER email;