-- --------------------------------------------------------
-- Host:                         localhost
-- Server version:               8.0.27 - MySQL Community Server - GPL
-- Server OS:                    Linux
-- HeidiSQL Version:             11.3.0.6295
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


-- Dumping database structure for go-todo
CREATE DATABASE IF NOT EXISTS `go-todo` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `go-todo`;

-- Dumping structure for table go-todo.todos
CREATE TABLE IF NOT EXISTS `todos` (
  `id` int NOT NULL AUTO_INCREMENT,
  `content` varchar(255) NOT NULL,
  `completed` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table go-todo.todos: ~0 rows (approximately)
/*!40000 ALTER TABLE `todos` DISABLE KEYS */;
INSERT IGNORE INTO `todos` (`id`, `content`, `completed`) VALUES
	(1, 'This is a new todo created from postman after authorization wiith NATS', 1);
/*!40000 ALTER TABLE `todos` ENABLE KEYS */;

-- Dumping structure for table go-todo.users
CREATE TABLE IF NOT EXISTS `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `password` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table go-todo.users: ~0 rows (approximately)
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT IGNORE INTO `users` (`id`, `username`, `password`) VALUES
	(1, 'tejas', '$2a$04$JBuKeCdiK.mkFfOriKAc6e/s5hH4IcH8xhLjzpt6VM4vY0HhfaBG.');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;

-- Dumping structure for table go-todo.users_todos
CREATE TABLE IF NOT EXISTS `users_todos` (
  `user_id` int NOT NULL,
  `todo_id` int NOT NULL,
  KEY `FK1_user_id` (`user_id`),
  KEY `FK2_todo_id` (`todo_id`),
  CONSTRAINT `FK1_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `FK2_todo_id` FOREIGN KEY (`todo_id`) REFERENCES `todos` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Dumping data for table go-todo.users_todos: ~0 rows (approximately)
/*!40000 ALTER TABLE `users_todos` DISABLE KEYS */;
INSERT IGNORE INTO `users_todos` (`user_id`, `todo_id`) VALUES
	(1, 1);
/*!40000 ALTER TABLE `users_todos` ENABLE KEYS */;

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
