# Load Balancer
resource "aws_security_group" "frontend_alb" {
  name_prefix = "${var.default_tags.project}-fe-alb"
  description = "Security Group for Frontend Load Balancer"
  vpc_id      = aws_vpc.main.id
}

resource "aws_security_group_rule" "frontend_alb_allow_80" {
  security_group_id = aws_security_group.frontend_alb.id
  type              = "ingress"
  protocol          = "tcp"
  from_port         = 80
  to_port           = 80
  cidr_blocks       = ["0.0.0.0/0"]
  ipv6_cidr_blocks  = ["::/0"]
  description       = "Allow HTTP traffic"
}

resource "aws_security_group_rule" "fronetned_alb_allow_outbound" {
  security_group_id = aws_security_group.frontend_alb.id
  type              = "egress"
  protocol          = "-1"
  from_port         = 0
  to_port           = 0
  cidr_blocks       = ["0.0.0.0/0"]
  ipv6_cidr_blocks  = ["::/0"]
  description       = "Allow outbound traffic"
}

# Frontend Servers
resource "aws_security_group" "frontend_server" {
  name_prefix = "${var.default_tags.project}-fe-server"
  description = "Security Group for Frontend Servers"
  vpc_id      = aws_vpc.main.id
}

resource "aws_security_group_rule" "frontend_server_allow_8080" {
  security_group_id = aws_security_group.frontend_server.id
  type              = "ingress"
  protocol          = "tcp"
  from_port         = 8080
  to_port           = 8080
  cidr_blocks       = ["0.0.0.0/0"]
  ipv6_cidr_blocks  = ["::/0"]
  description       = "Allow HTTP traffic"
}

resource "aws_security_group_rule" "frontend_server_allow_outbound" {
  security_group_id = aws_security_group.frontend_server.id
  type              = "egress"
  protocol          = "-1"
  from_port         = 0
  to_port           = 0
  cidr_blocks       = ["0.0.0.0/0"]
  ipv6_cidr_blocks  = ["::/0"]
  description       = "Allow outbound traffic"
}

# Dispatcher
resource "aws_security_group" "dispatcher" {
  name_prefix = "${var.default_tags.project}-dispatcher"
  description = "Security Group for Dispatcher"
  vpc_id      = aws_vpc.main.id
}

resource "aws_security_group_rule" "dispatcher_allow_8090" {
  security_group_id = aws_security_group.dispatcher.id
  type              = "ingress"
  protocol          = "tcp"
  from_port         = 8090
  to_port           = 8090
  cidr_blocks       = ["0.0.0.0/0"]
  ipv6_cidr_blocks  = ["::/0"]
  description       = "Allow HTTP traffic"
}

resource "aws_security_group_rule" "dispatcher_allow_outbound" {
  security_group_id = aws_security_group.dispatcher.id
  type              = "egress"
  protocol          = "-1"
  from_port         = 0
  to_port           = 0
  cidr_blocks       = ["0.0.0.0/0"]
  ipv6_cidr_blocks  = ["::/0"]
  description       = "Allow outbound traffic"
}
