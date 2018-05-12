create database vcentercollector;
GRANT USAGE ON `vcentercollector`.* to 'vcentercollectoruser'@'localhost' identified by 'vcentercollectorpass';
GRANT ALL PRIVILEGES ON `vcentercollector`.* to 'vcentercollectoruser'@'localhost' with grant option;
flush privileges;
