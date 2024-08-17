CREATE DATABASE  IF NOT EXISTS `assignment` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `assignment`;
-- MySQL dump 10.13  Distrib 8.0.36, for Win64 (x86_64)
--
-- Host: localhost    Database: assignment
-- ------------------------------------------------------
-- Server version	8.0.36

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `customers`
--

DROP TABLE IF EXISTS customers;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE customers (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  first_name longtext NOT NULL,
  last_name longtext NOT NULL,
  address longtext NOT NULL,
  phone longtext NOT NULL,
  balance double NOT NULL,
  deleted_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_customers_deleted_at (deleted_at)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `customers`
--

LOCK TABLES customers WRITE;
/*!40000 ALTER TABLE customers DISABLE KEYS */;
INSERT INTO customers VALUES (1,'ranj','azeez','hawler','34',790,NULL);
/*!40000 ALTER TABLE customers ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invoice_lines`
--

DROP TABLE IF EXISTS invoice_lines;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE invoice_lines (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  InvoiceId bigint unsigned DEFAULT NULL,
  item_id bigint unsigned DEFAULT NULL,
  quantity double DEFAULT '0',
  line_price double DEFAULT '0',
  deleted_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_invoice_lines_deleted_at (deleted_at),
  KEY fk_invoices_invoice_line (InvoiceId),
  CONSTRAINT fk_invoices_invoice_line FOREIGN KEY (InvoiceId) REFERENCES invoices (id)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invoice_lines`
--

LOCK TABLES invoice_lines WRITE;
/*!40000 ALTER TABLE invoice_lines DISABLE KEYS */;
INSERT INTO invoice_lines VALUES (1,1,1,10,10,NULL),(2,1,1,10,10,NULL);
/*!40000 ALTER TABLE invoice_lines ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invoices`
--

DROP TABLE IF EXISTS invoices;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE invoices (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  invoice_unique_id varchar(191) DEFAULT NULL,
  customer bigint unsigned DEFAULT NULL,
  invoice_total double DEFAULT NULL,
  invoice_date date DEFAULT NULL,
  deleted_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uni_invoices_invoice_unique_id (invoice_unique_id),
  KEY idx_invoices_deleted_at (deleted_at)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invoices`
--

LOCK TABLES invoices WRITE;
/*!40000 ALTER TABLE invoices DISABLE KEYS */;
INSERT INTO invoices VALUES (1,'2024-0001',1,210,'2024-08-17',NULL);
/*!40000 ALTER TABLE invoices ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `products`
--

DROP TABLE IF EXISTS products;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE products (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` longtext,
  barcode varchar(191) DEFAULT NULL,
  quantity_on_hand double DEFAULT NULL,
  price double DEFAULT NULL,
  supplier bigint unsigned DEFAULT NULL,
  product_image varchar(191) DEFAULT 'Null1',
  deleted_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uni_products_barcode (barcode),
  KEY idx_products_deleted_at (deleted_at),
  KEY fk_suppliers_product (supplier),
  CONSTRAINT fk_suppliers_product FOREIGN KEY (supplier) REFERENCES suppliers (id)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `products`
--

LOCK TABLES products WRITE;
/*!40000 ALTER TABLE products DISABLE KEYS */;
INSERT INTO products VALUES (1,'egg','e',170,10,1,'images\\a.jpg',NULL);
/*!40000 ALTER TABLE products ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `suppliers`
--

DROP TABLE IF EXISTS suppliers;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE suppliers (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(191) NOT NULL DEFAULT 'Null  size:40',
  phone varchar(30) NOT NULL,
  deleted_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uni_suppliers_phone (phone),
  KEY idx_suppliers_deleted_at (deleted_at)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `suppliers`
--

LOCK TABLES suppliers WRITE;
/*!40000 ALTER TABLE suppliers DISABLE KEYS */;
INSERT INTO suppliers VALUES (1,'harem','34342',NULL);
/*!40000 ALTER TABLE suppliers ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS users;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE users (
  id bigint unsigned NOT NULL AUTO_INCREMENT,
  user_name varchar(191) NOT NULL,
  `password` longtext NOT NULL,
  `role` varchar(191) NOT NULL DEFAULT 'member',
  deleted_at datetime(3) DEFAULT NULL,
  PRIMARY KEY (id),
  UNIQUE KEY uni_users_user_name (user_name),
  KEY idx_users_delete_at (deleted_at)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES users WRITE;
/*!40000 ALTER TABLE users DISABLE KEYS */;
INSERT INTO users VALUES (1,'admin','$2a$10$4omNO/zT63HX3RfCDKZ2jeJEcPdwzyijp3aXDj3rLgQnxvcf5DOhG','Admin',NULL),(2,'ranj','$2a$10$heJtrEzXgweZcgSayfThj.Cqd80MeDSld5VC51gQNnfbKwye2kEuK','member',NULL);
/*!40000 ALTER TABLE users ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-08-17 21:21:36
