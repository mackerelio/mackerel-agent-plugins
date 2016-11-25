FROM php:5.5-apache

ENV DEBIAN_FRONTEND noninteractive

RUN docker-php-ext-install opcache
RUN mkdir /var/www/html/mackerel
ADD php-opcache.php /var/www/html/mackerel/
