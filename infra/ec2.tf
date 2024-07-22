resource "aws_instance" "dispatcher" {
  ami                         = "ami-04a81a99f5ec58529"
  instance_type               = "t2.micro"
  subnet_id                   = aws_subnet.public[0].id
  associate_public_ip_address = true
  key_name                    = var.dispatcher_key_pair
  vpc_security_group_ids      = [aws_security_group.dispatcher.id]
  monitoring                  = true

  user_data = base64encode(file("${path.module}/scripts/dispatcher.sh"))

  tags = {
    "Name" = "${var.default_tags.project}-dispatcher"
  }
}

resource "aws_instance" "frontend_server" {
  count = var.frontend_server_count

  ami                         = "ami-04a81a99f5ec58529"
  instance_type               = "t2.micro"
  subnet_id                   = aws_subnet.public[count.index % var.public_subnet_count].id
  associate_public_ip_address = true
  key_name                    = var.frontend_server_key_pair
  vpc_security_group_ids      = [aws_security_group.frontend_server.id]
  monitoring                  = true

  user_data = base64encode(templatefile("${path.module}/scripts/fe_server.sh", {
    DISPATCHER_PUBLIC_DNS = data.aws_instance.dispatcher.public_dns
  }))

  tags = {
    "Name" = "${var.default_tags.project}-fe-server"
  }

  depends_on = [aws_instance.dispatcher]
}
