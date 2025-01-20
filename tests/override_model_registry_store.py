from mlflow.store.model_registry.sqlalchemy_store import SqlAlchemyStore

from mlflow_go_backend.store.model_registry import ModelRegistryStore

SqlAlchemyStore = ModelRegistryStore(SqlAlchemyStore)
