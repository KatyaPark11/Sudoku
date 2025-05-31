from db.models import UserModel
from schemas.users import UserCreate
from sqlalchemy.orm import Session


def get_users(db: Session):
    return db.query(UserModel).all()


def get_user_by_id(db: Session, user_id: int):
    return db.query(UserModel).filter(UserModel.id == user_id).first()


def get_user_by_name(db: Session, username: int):
    return db.query(UserModel).filter(UserModel.username == username).first()


def create_user(db: Session, user: UserCreate):
    db_user = UserModel(**user.dict())
    db.add(db_user)
    db.commit()
    db.refresh(db_user)
    return db_user


def update_user(db: Session, db_user: UserModel, user: UserCreate):
    for key, value in user.dict().items():
        setattr(db_user, key, value)
    db.commit()
    return db_user


def delete_user(db: Session, db_user: UserModel):
    db.delete(db_user)
    db.commit()
    return db_user