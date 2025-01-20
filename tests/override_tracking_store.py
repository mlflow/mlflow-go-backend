from mlflow.store.tracking.sqlalchemy_store import SqlAlchemyStore

from mlflow_go_backend.store.tracking import TrackingStore

SqlAlchemyStore = TrackingStore(SqlAlchemyStore)
