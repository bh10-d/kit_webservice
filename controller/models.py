from flask_sqlalchemy import SQLAlchemy
import uuid

db = SQLAlchemy()

class Runner(db.Model):
    __tablename__ = 'runners'
    id = db.Column(db.String, primary_key=True, unique=True)  # runner_id
    name = db.Column(db.String, nullable=False)
    ip = db.Column(db.String, nullable=False)
    tags = db.Column(db.String, nullable=False)
    # capacity = db.Column(db.Integer, default=1)
    # runner_token = db.Column(db.String, unique=True, nullable=False)
    created_at = db.Column(db.DateTime, default=db.func.now())  # 👈 Thời gian tạo

    def to_dict(self):
        return {
            "id": self.id,
            "name": self.name,
            "ip": self.ip,
            # "tags": self.tags,
            # "capacity": self.capacity,
            # "runner_token": self.runner_token
        }

from datetime import datetime, timezone

class Job(db.Model):
    __tablename__ = 'job'

    id = db.Column(db.String, primary_key=True, default=lambda: str(uuid.uuid4()))
    # id = db.Column(db.String, primary_key=True)
    runner_id = db.Column(db.String, db.ForeignKey('runners.id'), nullable=False)
    msg_id = db.Column(db.String, nullable=False)
    status = db.Column(db.String, nullable=False, default='pending')
    request_payload = db.Column(db.Text, nullable=False)
    response_payload = db.Column(db.Text)
    error_message = db.Column(db.Text)
    timeout = db.Column(db.Boolean, default=False)
    created_at = db.Column(db.DateTime, default=lambda: datetime.now(timezone.utc))
    updated_at = db.Column(db.DateTime, default=lambda: datetime.now(timezone.utc), onupdate=lambda: datetime.now(timezone.utc))
